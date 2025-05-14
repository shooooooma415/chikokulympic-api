package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SigninRequest struct {
	AuthID entity.AuthID `json:"auth_id" example:"auth456"`
}

type SigninResponse struct {
	UserID entity.UserID `json:"user_id" example:"user123"`
}

type Signin struct {
	userRepo repository.UserRepository
}

func NewSignin(userRepo repository.UserRepository) *Signin {
	return &Signin{
		userRepo: userRepo,
	}
}

// @Summary signin user
// @Description signin user from auth_id
// @Tags users
// @Accept json
// @Produce json
// @Param request body SigninRequest true "request"
// @Success 200 {object} SigninResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /users/signin [post]
func (s *Signin) Handler(c echo.Context) error {
	req := new(SigninRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	if req.AuthID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("AuthIDは必須です"))
	}

	authID := entity.AuthID(req.AuthID)
	user, err := usecase.NewAuthenticateUserUseCase(s.userRepo, authID).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	response := SigninResponse{
		UserID: user.UserID,
	}

	return c.JSON(http.StatusOK, response)
}
