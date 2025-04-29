package user

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type userUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{userRepo: r}
}

func (u *userUsecaseImpl) RegisterUser(user *entity.User) (*entity.User, error) {
	return u.userRepo.CreateUser(user)
}

func (u *userUsecaseImpl) AuthenticateUser(authID entity.AuthID) (*entity.User, error) {
	return u.userRepo.FindUserByAuthID(authID)
}

func (u *userUsecaseImpl) UpdateUserName(user *entity.User) (*entity.User, error) {
	return u.userRepo.UpdateUser(user)
}