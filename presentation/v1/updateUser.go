package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)


type UpdateUserRequest struct {
	UserName entity.UserName `json:"user_name" example:"user_name"`
	UserIcon entity.UserIcon `json:"user_icon" example:"https://example.com/icon.png"`
}


type UpdateUserResponse struct {
	UserID entity.UserID `json:"user_id" example:"user123"`
}

type UpdateUser struct {
	userRepo repository.UserRepository
}

func NewUpdateUser(userRepo repository.UserRepository) *UpdateUser {
	return &UpdateUser{
		userRepo: userRepo,
	}
}


// @Summary update user
// @Description update user information(icon, name)
// @Tags users
// @Accept json
// @Produce json
// @Param request body UpdateUserRequest true "request"
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
