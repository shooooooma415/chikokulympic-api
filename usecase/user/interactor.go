package user

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(r repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: r}
}

func (u *userUsecase) RegisterUser(user *entity.User) (*entity.User, error) {
	return u.userRepo.CreateUser(user)
}

func (u *userUsecase) AuthenticateUser(authID entity.AuthID) (*entity.User, error) {
	return u.userRepo.FindUserByAuthID(authID)
}

func (u *userUsecase) UpdateUserName(user *entity.User) (*entity.User, error) {
	return u.userRepo.UpdateUser(user)
}