package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// LeaveGroupRequest はグループ退会リクエストの構造体
// @Description グループ退会リクエスト
type LeaveGroupRequest struct {
	UserID entity.UserID `json:"user_id" validate:"required" example:"user123"`
}

// LeaveGroup はグループ退会ハンドラの構造体
type LeaveGroup struct {
	groupRepo repository.GroupRepository
}

// NewLeaveGroup は新しいLeaveGroupハンドラを作成する
func NewLeaveGroup(groupRepo repository.GroupRepository) *LeaveGroup {
	return &LeaveGroup{
		groupRepo: groupRepo,
	}
}

// Handler はグループ退会APIのハンドラ
// @Summary グループから退会する
// @Description ユーザーが指定したグループから退会する
// @Tags groups
// @Accept json
// @Produce json
// @Param group_id path string true "グループID"
// @Param request body LeaveGroupRequest true "グループ退会リクエスト"
// @Success 200 {object} nil
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /groups/{group_id}/leave [post]
func (l *LeaveGroup) Handler(c echo.Context) error {
	// パスパラメータからgroup_idを取得
	groupIDParam := c.Param("group_id")
	if groupIDParam == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("グループIDは必須です"))
	}

	// リクエストボディからユーザーIDを取得
	req := new(LeaveGroupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	// バリデーション
	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("ユーザーIDは必須です"))
	}

	groupID := entity.GroupID(groupIDParam)
	userID := entity.UserID(req.UserID)

	err := usecase.NewLeaveGroupUseCase(l.groupRepo, userID, groupID).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
