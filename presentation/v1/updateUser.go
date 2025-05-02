package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UpdateUserRequest struct {
	UserName string `json:"user_name"`
	UserIcon string `json:"user_icon"`
}

type UpdateUserResponse struct {
	UserID string `json:"user_id"`
}

type UpdateUser struct {
	usecase.UpdateUserUseCase
}

func NewUpdateUser(usecase usecase.UpdateUserUseCase) *UpdateUser {
	return &UpdateUser{
		UpdateUserUseCase: usecase,
	}
}
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

	updatedUser, err := u.Execute(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	response := UpdateUserResponse{
		UserID: string(updatedUser.UserID),
	}

	return c.JSON(http.StatusOK, response)
}
