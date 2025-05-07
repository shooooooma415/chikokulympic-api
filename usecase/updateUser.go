package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type UpdateUserUseCase interface {
	Execute() (*entity.User, error)
}

type UpdateUserUseCaseImpl struct {
	userRepo repository.UserRepository
	user     *entity.User
}

func NewUpdateUserUseCase(userRepo repository.UserRepository, user *entity.User) *UpdateUserUseCaseImpl {
	return &UpdateUserUseCaseImpl{
		userRepo: userRepo,
		user:     user,
	}
}

func (uc *UpdateUserUseCaseImpl) Execute() (*entity.User, error) {
	return uc.userRepo.UpdateUser(uc.user)
}
