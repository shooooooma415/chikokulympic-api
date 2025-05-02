package v1

import (
	presentationV1 "chikokulympic-api/presentation/v1"
	"chikokulympic-api/usecase"

	"github.com/labstack/echo/v4"
)

type AuthServer struct {
	signup     *presentationV1.Signup
	signin     *presentationV1.Signin
	updateUser *presentationV1.UpdateUser
}

func NewAuthServer(
	registerUserUseCase usecase.RegisterUserUseCase,
	authenticateUserUseCase usecase.AuthenticateUserUseCase,
	updateUserUseCase usecase.UpdateUserUseCase,
) *AuthServer {
	return &AuthServer{
		signup:     presentationV1.NewSignup(registerUserUseCase),
		signin:     presentationV1.NewSignin(authenticateUserUseCase),
		updateUser: presentationV1.NewUpdateUser(updateUserUseCase),
	}
}

func (s *AuthServer) RegisterRoutes(e *echo.Echo) {
	authGroup := e.Group("/auth")

	// サインアップ用エンドポイント
	authGroup.POST("/signup", s.signup.Handler)

	// サインイン用エンドポイント
	authGroup.POST("/signin", s.signin.Handler)

	// ユーザー情報更新用エンドポイント
	authGroup.PUT("", s.updateUser.Handler)
}
