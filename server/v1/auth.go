package v1

import (
	"chikokulympic-api/domain/repository"
	presentationV1 "chikokulympic-api/presentation/v1"

	"github.com/labstack/echo/v4"
)

type AuthServer struct {
	signup     *presentationV1.Signup
	signin     *presentationV1.Signin
	updateUser *presentationV1.UpdateUser
}

func NewAuthServer(userRepo repository.UserRepository) *AuthServer {
	return &AuthServer{
		signup:     presentationV1.NewSignup(userRepo),
		signin:     presentationV1.NewSignin(userRepo),
		updateUser: presentationV1.NewUpdateUser(userRepo),
	}
}

func (s *AuthServer) RegisterRoutes(e *echo.Echo) {
	authGroup := e.Group("/auth")

	authGroup.POST("/signup", s.signup.Handler)

	authGroup.POST("/signin", s.signin.Handler)

	authGroup.PUT("", s.updateUser.Handler)
}
