package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SignupRequest struct {
	FCMToken entity.FCMToken `json:"token" validate:"required" example:"fcm_token"`
	UserName entity.UserName `json:"user_name" validate:"required" example:"user_name"`
	AuthID   entity.AuthID   `json:"auth_id" validate:"required" example:"auth_id"`
	UserIcon entity.UserIcon `json:"user_icon" example:"user_icon"`
}

type SignupResponse struct {
	UserID entity.UserID `json:"user_id" example:"user123"`
}

type Signup struct {
	userRepo repository.UserRepository
}

func NewSignup(userRepo repository.UserRepository) *Signup {
	return &Signup{
		userRepo: userRepo,
	}
}

// @Summary subscribe user
// @Description subscribe user to the service
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignupRequest true "request"
// @Success 201 {object} SignupResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /users/signup [post]
func (s *Signup) Handler(c echo.Context) error {
	req := new(SignupRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	if req.UserName == "" || req.AuthID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("ユーザー名とAuthIDは必須です"))
	}

	user := &entity.User{
		AuthID:   entity.AuthID(req.AuthID),
		UserName: entity.UserName(req.UserName),
		FCMToken: entity.FCMToken(req.FCMToken),
		UserIcon: entity.UserIcon(req.UserIcon),
	}

	registeredUser, err := usecase.NewRegisterUserUseCase(s.userRepo, user).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	response := SignupResponse{
		UserID: registeredUser.UserID,
	}

	return c.JSON(http.StatusCreated, response)
}
