package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PostVoteRequest struct {
	UserID entity.UserID `json:"user_id" validate:"required" example:"user123"`
	Option entity.Vote   `json:"option" validate:"required" example:"参加"`
}

type PostVote struct {
	eventRepo repository.EventRepository
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

func NewPostVote(eventRepo repository.EventRepository, groupRepo repository.GroupRepository, userRepo repository.UserRepository) *PostVote {
	return &PostVote{
		eventRepo: eventRepo,
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// @Summary post vote
// @Description post a vote for an event
// @Tags events
// @Accept json
// @Produce json
// @Param event_id path string true "Event ID"
// @Param request body PostVoteRequest true "request"
// @Success 200
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 403 {object} middleware.ErrorResponse
// @Failure 404 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /events/{event_id}/votes [post]
func (p *PostVote) Handler(c echo.Context) error {
	eventIDStr := c.Param("event_id")
	if eventIDStr == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("イベントIDは必須です"))
	}

	eventID := entity.EventID(eventIDStr)

	req := new(PostVoteRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("ユーザーIDは必須です"))
	}
	if req.Option == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("投票オプションは必須です"))
	}

	userIDInt, err := strconv.Atoi(string(req.UserID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("無効なユーザーIDです"))
	}
	userID := entity.UserID(strconv.Itoa(userIDInt))

	_, err = usecase.NewPostParticipationUseCase(p.eventRepo, p.groupRepo, &userID, &eventID, &req.Option).Execute()
	if err != nil {
		if err.Error() == "event not found" {
			return c.JSON(http.StatusNotFound, middleware.NewErrorResponse("イベントが見つかりません"))
		}
		if err.Error() == "not a group member. cannot vote" {
			return c.JSON(http.StatusForbidden, middleware.NewErrorResponse("このイベントに投票する権限がありません。グループに所属しているか確認してください"))
		}
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
