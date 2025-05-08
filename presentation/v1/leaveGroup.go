package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LeaveGroupRequest struct {
	GroupName entity.GroupName `json:"group_name" validate:"required"`
	UserID    entity.UserID    `json:"user_id" validate:"required"`
}

type LeaveGroup struct {
	groupRepo repository.GroupRepository
}

func NewLeaveGroup(groupRepo repository.GroupRepository) *LeaveGroup {
	return &LeaveGroup{
		groupRepo: groupRepo,
	}
}

func (l *LeaveGroup) Handler(c echo.Context) error {
	req := new(LeaveGroupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	// バリデーション
	if req.GroupName == "" || req.UserID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("グループ名とユーザーIDは必須です"))
	}

	group := &entity.Group{
		GroupName: req.GroupName,
	}

	userID := entity.UserID(req.UserID)

	err := usecase.NewLeaveGroupUseCase(l.groupRepo, userID, *group).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
