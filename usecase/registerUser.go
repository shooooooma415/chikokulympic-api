package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"

	"github.com/google/uuid"
)

type RegisterUserUseCase interface {
	Execute(user *entity.User) (*entity.User, error)
}

type RegisterUserUseCaseImpl struct {
	userRepo repository.UserRepository
}

func NewRegisterUserUseCase(userRepo repository.UserRepository) *RegisterUserUseCaseImpl {
	return &RegisterUserUseCaseImpl{
		userRepo: userRepo,
	}
}

func (uc *RegisterUserUseCaseImpl) Execute(user *entity.User) (*entity.User, error) {
	userID := uuid.New().String()
	user.UserID = entity.UserID(userID)

	return uc.userRepo.CreateUser(user)
}
