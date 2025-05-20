package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
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
	return uc.userRepo.CreateUser(*uc.user)
}
