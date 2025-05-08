package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type JoinGroupRequest struct {
	GroupName     entity.GroupName     `json:"group_name" validate:"required"`
	GroupPassword entity.GroupPassword `json:"group_password" validate:"required"`
	UserID        entity.UserID        `json:"user_id" validate:"required"`
}

type JoinGroupResponse struct {
	GroupID entity.GroupID `json:"group_id"`
}

type JoinGroup struct {
	userRepo  repository.UserRepository
	groupRepo repository.GroupRepository
}

func NewJoinGroup(userRepo repository.UserRepository, groupRepo repository.GroupRepository) *JoinGroup {
	return &JoinGroup{
		userRepo:  userRepo,
		groupRepo: groupRepo,
	}
}
func (j *JoinGroup) Handler(c echo.Context) error {
	req := new(JoinGroupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	if req.GroupName == "" || req.GroupPassword == "" || req.UserID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("グループ名、パスワード、ユーザーIDは必須です"))
	}

	group := &entity.Group{
		GroupName:     req.GroupName,
		GroupPassword: req.GroupPassword,
	}

	userID := entity.UserID(req.UserID)

	groupID, err := usecase.NewJoinGroupUseCase(j.groupRepo, j.userRepo, userID, *group).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	response := JoinGroupResponse{
		GroupID: *groupID,
	}

	return c.JSON(http.StatusOK, response)
}
