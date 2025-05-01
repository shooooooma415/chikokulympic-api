package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type AuthenticateUserUseCase interface {
	Execute(authID entity.AuthID) (*entity.User, error)
}

type AuthenticateUserUseCaseImpl struct {
	userRepo repository.UserRepository
}

func NewAuthenticateUserUseCase(userRepo repository.UserRepository) *AuthenticateUserUseCaseImpl {
	return &AuthenticateUserUseCaseImpl{
		userRepo: userRepo,
	}
}

func (uc *AuthenticateUserUseCaseImpl) Execute(authID entity.AuthID) (*entity.User, error) {
	return uc.userRepo.FindUserByAuthID(authID)
}
