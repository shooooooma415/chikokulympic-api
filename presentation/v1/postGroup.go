package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PostGroupRequest struct {
	GroupName        entity.GroupName        `json:"group_name" validate:"required"`
	GroupPassword    entity.GroupPassword    `json:"group_password" validate:"required"`
	ManagerID        entity.UserID           `json:"manager_id" validate:"required"`
	GroupDescription entity.GroupDescription `json:"group_description" validate:"required"`
}

type PostGroupResponse struct {
	GroupID entity.GroupID `json:"group_id"`
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

func (p *PostGroup) Handler(c echo.Context) error {
	req := new(PostGroupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	if req.GroupName == "" || req.GroupPassword == "" || req.ManagerID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("グループ名、パスワード、作成者IDは必須です"))
	}

	groupID := uuid.New().String()
	group := &entity.Group{
		GroupID:          entity.GroupID(groupID),
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
