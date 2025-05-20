package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"

	"github.com/google/uuid"
)

type RegisterUserUseCase interface {
	Execute() (*entity.User, error)
}

type RegisterUserUseCaseImpl struct {
	userRepo repository.UserRepository
	user     *entity.User
}

func NewRegisterUserUseCase(userRepo repository.UserRepository, user *entity.User) *RegisterUserUseCaseImpl {
	return &RegisterUserUseCaseImpl{
		userRepo: userRepo,
		user:     user,
	}
}

func (uc *RegisterUserUseCaseImpl) Execute() (*entity.User, error) {
	userID := uuid.New().String()
	uc.user.UserID = entity.UserID(userID)

	return uc.userRepo.CreateUser(*uc.user)
}
