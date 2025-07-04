package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostGroupRequest struct {
	GroupName        entity.GroupName        `json:"group_name" validate:"required" example:"group_name"`
	GroupPassword    entity.GroupPassword    `json:"password" validate:"required" example:"password"`
	ManagerID        entity.UserID           `json:"manager_id" validate:"required" example:"user_id"`
	GroupDescription entity.GroupDescription `json:"description" validate:"required" example:"description"`
}

type PostGroupResponse struct {
	GroupID entity.GroupID `json:"group_id" example:"group123"`
}

type PostGroup struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

func NewPostGroup(groupRepo repository.GroupRepository, userRepo repository.UserRepository) *PostGroup {
	return &PostGroup{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// @Summary create group
// @Description create a new group
// @Tags groups
// @Accept json
// @Produce json
// @Param request body PostGroupRequest true "request"
// @Success 201 {object} PostGroupResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /groups [post]
func (p *PostGroup) Handler(c echo.Context) error {
	req := new(PostGroupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	if req.GroupName == "" || req.GroupPassword == "" || req.ManagerID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("グループ名、パスワード、作成者IDは必須です"))
	}

	group := &entity.Group{
		GroupName:        req.GroupName,
		GroupPassword:    req.GroupPassword,
		GroupManagerID:   req.ManagerID,
		GroupDescription: req.GroupDescription,
		GroupMembers:     entity.GroupMembers{},
		GroupEvents:      entity.GroupEvents{},
	}

	createdGroup, err := usecase.NewCreateGroupUseCase(p.groupRepo, p.userRepo, group).Execute()
	if err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	response := PostGroupResponse{
		GroupID: createdGroup.GroupID,
	}

	return c.JSON(http.StatusCreated, response)
}
