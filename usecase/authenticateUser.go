package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type AuthenticateUserUseCase interface {
	Execute() (*entity.User, error)
}

type AuthenticateUserUseCaseImpl struct {
	userRepo repository.UserRepository
	authID   entity.AuthID
}

func NewAuthenticateUserUseCase(userRepo repository.UserRepository, authID entity.AuthID) *AuthenticateUserUseCaseImpl {
	return &AuthenticateUserUseCaseImpl{
		userRepo: userRepo,
		authID:   authID,
	}
}

func (uc *AuthenticateUserUseCaseImpl) Execute() (*entity.User, error) {
	return uc.userRepo.FindUserByAuthID(uc.authID)
}
