package v1

import (
	"chikokulympic-api/domain/repository"
	presentationV1 "chikokulympic-api/presentation/v1"

	"github.com/labstack/echo/v4"
)

type UserServer struct {
	signup        *presentationV1.Signup
	signin        *presentationV1.Signin
	updateUser    *presentationV1.UpdateUser
	getUserGroups *presentationV1.GetUserGroups
}

func NewUserServer(userRepo repository.UserRepository, groupRepo repository.GroupRepository) *UserServer {
	return &UserServer{
		signup:        presentationV1.NewSignup(userRepo),
		signin:        presentationV1.NewSignin(userRepo),
		updateUser:    presentationV1.NewUpdateUser(userRepo),
		getUserGroups: presentationV1.NewGetUserGroups(groupRepo),
	}
}

func (s *UserServer) RegisterRoutes(e *echo.Echo) {
	authGroup := e.Group("/users")

	authGroup.POST("/signup", s.signup.Handler)

	authGroup.POST("/signin", s.signin.Handler)

	authGroup.PUT("", s.updateUser.Handler)

	authGroup.GET("/:user_id/groups", s.getUserGroups.Handler)
}
