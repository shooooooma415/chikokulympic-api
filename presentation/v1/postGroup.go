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

// PostGroupRequest はグループ作成リクエストの構造体
// @Description グループ作成リクエスト
type PostGroupRequest struct {
	GroupName        entity.GroupName        `json:"group_name" example:"テストグループ"`
	GroupPassword    entity.GroupPassword    `json:"group_password" example:"password123"`
	ManagerID        entity.UserID           `json:"manager_id" example:"user123"`
	GroupDescription entity.GroupDescription `json:"group_description" example:"これはテストグループです"`
}

// PostGroupResponse はグループ作成レスポンスの構造体
// @Description グループ作成レスポンス
type PostGroupResponse struct {
	GroupID entity.GroupID `json:"group_id" example:"group123"`
}

// PostGroup はグループ作成ハンドラの構造体
type PostGroup struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

// NewPostGroup は新しいPostGroupハンドラを作成する
func NewPostGroup(groupRepo repository.GroupRepository, userRepo repository.UserRepository) *PostGroup {
	return &PostGroup{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

// Handler はグループ作成APIのハンドラ
// @Summary 新しいグループを作成する
// @Description 新しいグループを作成する
// @Tags groups
// @Accept json
// @Produce json
// @Param request body PostGroupRequest true "グループ作成リクエスト"
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
