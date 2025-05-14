package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// UpdateUserRequest はユーザー情報更新リクエストの構造体
// @Description ユーザー情報更新リクエスト
type UpdateUserRequest struct {
	UserName entity.UserName `json:"user_name" example:"山田太郎"`
	UserIcon entity.UserIcon `json:"user_icon" example:"https://example.com/icon.png"`
}

// UpdateUserResponse はユーザー情報更新レスポンスの構造体
// @Description ユーザー情報更新レスポンス
type UpdateUserResponse struct {
	UserID entity.UserID `json:"user_id" example:"user123"`
}

// UpdateUser はユーザー情報更新ハンドラの構造体
type UpdateUser struct {
	userRepo repository.UserRepository
}

// NewUpdateUser は新しいUpdateUserハンドラを作成する
func NewUpdateUser(userRepo repository.UserRepository) *UpdateUser {
	return &UpdateUser{
		userRepo: userRepo,
	}
}

// Handler はユーザー情報更新APIのハンドラ
// @Summary ユーザー情報を更新する
// @Description ユーザーの名前やアイコンなどの情報を更新する
// @Tags users
// @Accept json
// @Produce json
// @Param request body UpdateUserRequest true "ユーザー情報更新リクエスト"
// @Success 200 {object} UpdateUserResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /users [put]
func (u *UpdateUser) Handler(c echo.Context) error {
	req := new(UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	if req.UserName == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("ユーザー名は必須です"))
	}

	user := &entity.User{
		UserName: entity.UserName(req.UserName),
		UserIcon: entity.UserIcon(req.UserIcon),
	}

	updatedUser, err := usecase.NewUpdateUserUseCase(u.userRepo, user).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	response := UpdateUserResponse{
		UserID: updatedUser.UserID,
	}

	return c.JSON(http.StatusOK, response)
}
