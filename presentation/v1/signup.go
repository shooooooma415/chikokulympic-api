package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

// SignupRequest はサインアップリクエストの構造体
// @Description サインアップリクエスト
type SignupRequest struct {
	FCMToken entity.FCMToken `json:"token" example:"fcm-token-123456"`
	UserName entity.UserName `json:"user_name" example:"山田太郎"`
	AuthID   entity.AuthID   `json:"auth_id" example:"auth456"`
	UserIcon entity.UserIcon `json:"user_icon" example:"https://example.com/icon.png"`
}

// SignupResponse はサインアップレスポンスの構造体
// @Description サインアップレスポンス
type SignupResponse struct {
	UserID entity.UserID `json:"user_id" example:"user123"`
}

// Signup はサインアップハンドラの構造体
type Signup struct {
	userRepo repository.UserRepository
}

// NewSignup は新しいSignupハンドラを作成する
func NewSignup(userRepo repository.UserRepository) *Signup {
	return &Signup{
		userRepo: userRepo,
	}
}

// Handler はサインアップAPIのハンドラ
// @Summary 新規ユーザー登録を行う
// @Description 新規ユーザーの登録（サインアップ）を行う
// @Tags users
// @Accept json
// @Produce json
// @Param request body SignupRequest true "サインアップリクエスト"
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
