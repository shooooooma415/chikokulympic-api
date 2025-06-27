package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type GetEventBoard struct {
	groupRepo repository.GroupRepository
	eventRepo repository.EventRepository
	userRepo  repository.UserRepository
}

func NewGetEventBoard(groupRepo repository.GroupRepository, eventRepo repository.EventRepository, userRepo repository.UserRepository) *GetEventBoard {
	return &GetEventBoard{
		groupRepo: groupRepo,
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

// @Summary get event board
// @Description get events by group IDs for event board
// @Tags events
// @Accept json
// @Produce json
// @Param group_ids query string true "Comma-separated list of group IDs"
// @Success 200 {object} usecase.FetchEventBoardResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /events/board [get]
func (g *GetEventBoard) Handler(c echo.Context) error {
	groupIDsParam := c.QueryParam("group_ids")
	if groupIDsParam == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("グループIDは必須です"))
	}

	groupIDStrings := strings.Split(groupIDsParam, ",")
	var groupIDs []entity.GroupID

	for _, idStr := range groupIDStrings {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		// グループIDを検証（空文字でないことを確認）
		if idStr == "" {
			return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(fmt.Sprintf("無効なグループID: %s", idStr)))
		}

		groupIDs = append(groupIDs, entity.GroupID(idStr))
	}

	if len(groupIDs) == 0 {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("有効なグループIDが指定されていません"))
	}

	result, err := usecase.NewFetchEventBoardUseCase(g.groupRepo, g.eventRepo, g.userRepo, groupIDs).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(fmt.Sprintf("イベントボードの取得中にエラーが発生しました: %v", err)))
	}

	return c.JSON(http.StatusOK, result)
}
