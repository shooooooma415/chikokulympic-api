package user

import "chikokulympic-api/domain/entity"

type UserUsecase interface {
	RegisterUser(user *entity.User) (*entity.User, error)
	AuthenticateUser(authID entity.AuthID) (*entity.User, error)
	UpdateUserName(user *entity.User) (*entity.User, error)
}
