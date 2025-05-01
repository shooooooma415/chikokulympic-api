package user

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type UpdateUserNameUseCase interface {
	Execute(user *entity.User) (*entity.User, error)
}

type UpdateUserNameUseCaseImpl struct {
	userRepo repository.UserRepository
}

func NewUpdateUserNameUseCase(userRepo repository.UserRepository) *UpdateUserNameUseCaseImpl {
	return &UpdateUserNameUseCaseImpl{
		userRepo: userRepo,
	}
}

func (uc *UpdateUserNameUseCaseImpl) Execute(user *entity.User) (*entity.User, error) {
	return uc.userRepo.UpdateUser(user)
}
