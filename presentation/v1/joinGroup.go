package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// JoinGroupRequest はグループ参加リクエストの構造体
// @Description グループ参加リクエスト
type JoinGroupRequest struct {
	GroupName     entity.GroupName     `json:"group_name" validate:"required" example:"テストグループ"`
	GroupPassword entity.GroupPassword `json:"group_password" validate:"required" example:"password123"`
	UserID        entity.UserID        `json:"user_id" validate:"required" example:"user123"`
}

// JoinGroupResponse はグループ参加レスポンスの構造体
// @Description グループ参加レスポンス
type JoinGroupResponse struct {
	GroupID entity.GroupID `json:"group_id" example:"group123"`
}

// JoinGroup はグループ参加ハンドラの構造体
type JoinGroup struct {
	userRepo  repository.UserRepository
	groupRepo repository.GroupRepository
}

// NewJoinGroup は新しいJoinGroupハンドラを作成する
func NewJoinGroup(userRepo repository.UserRepository, groupRepo repository.GroupRepository) *JoinGroup {
	return &JoinGroup{
		userRepo:  userRepo,
		groupRepo: groupRepo,
	}
}

// Handler はグループ参加APIのハンドラ
// @Summary グループに参加する
// @Description ユーザーが指定したグループに参加する
// @Tags groups
// @Accept json
// @Produce json
// @Param request body JoinGroupRequest true "グループ参加リクエスト"
// @Success 200 {object} JoinGroupResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /groups/join [post]
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
