package user

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
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
	return uc.userRepo.CreateUser(user)
}
