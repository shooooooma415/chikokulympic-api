package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SignupRequest struct {
	FCMToken string `json:"token"`
	UserName string `json:"user_name"`
	AuthID   string `json:"auth_id"`
}

type SignupResponse struct {
	UserID string `json:"user_id"`
}

type Signup struct {
	usecase.RegisterUserUseCase
}

func NewSignup(usecase usecase.RegisterUserUseCase) *Signup {
	return &Signup{
		RegisterUserUseCase: usecase,
	}
}

func (s *Signup) Handle(c echo.Context) error {
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
	}

	registeredUser, err := s.Execute(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	response := SignupResponse{
		UserID: string(registeredUser.UserID),
	}

	return c.JSON(http.StatusCreated, response)
}
