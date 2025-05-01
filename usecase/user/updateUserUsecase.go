package user

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type UpdateUserUseCase interface {
	Execute(user *entity.User) (*entity.User, error)
}

type UpdateUserUseCaseImpl struct {
	userRepo repository.UserRepository
}

func NewUpdateUserUseCase(userRepo repository.UserRepository) *UpdateUserUseCaseImpl {
	return &UpdateUserUseCaseImpl{
		userRepo: userRepo,
	}
}

func (uc *UpdateUserUseCaseImpl) Execute(user *entity.User) (*entity.User, error) {
	return uc.userRepo.UpdateUser(user)
}
