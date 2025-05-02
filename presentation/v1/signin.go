package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type SigninRequest struct {
	AuthID string `json:"auth_id"`
}

type SigninResponse struct {
	UserID string `json:"user_id"`
}

type Signin struct {
	usecase.AuthenticateUserUseCase
}

func NewSignin(usecase usecase.AuthenticateUserUseCase) *Signin {
	return &Signin{
		AuthenticateUserUseCase: usecase,
	}
}
func (s *Signin) Handle(c echo.Context) error {
	req := new(SigninRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse(err.Error()))
	}

	if req.AuthID == "" {
		return c.JSON(http.StatusBadRequest, middleware.NewErrorResponse("AuthIDは必須です"))
	}

	user, err := s.Execute(entity.AuthID(req.AuthID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	response := SigninResponse{
		UserID: string(user.UserID),
	}

	return c.JSON(http.StatusOK, response)
}
