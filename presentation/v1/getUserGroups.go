package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetUserGroups はユーザーのグループ一覧取得ハンドラの構造体
type GetUserGroups struct {
	groupRepo repository.GroupRepository
}

// NewGetUserGroups は新しいGetUserGroupsハンドラを作成する
func NewGetUserGroups(groupRepo repository.GroupRepository) *GetUserGroups {
	return &GetUserGroups{
		groupRepo: groupRepo,
	}
}

// Handler はユーザーのグループ一覧取得APIのハンドラ
// @Summary ユーザーが所属するグループ一覧を取得する
// @Description 指定したユーザーIDのユーザーが所属するグループの一覧を取得する
// @Tags groups
// @Accept json
// @Produce json
// @Param user_id path string true "ユーザーID"
// @Success 200 {array} usecase.UserGroup
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /users/{user_id}/groups [get]
func (g *GetUserGroups) Handler(c echo.Context) error {
	userIDParam := c.Param("user_id")

	if userIDParam == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("ユーザーIDは必須です"))
	}

	userID := entity.UserID(userIDParam)

	result, err := usecase.NewFetchUserGroupsUseCase(g.groupRepo, userID).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, result)
}
