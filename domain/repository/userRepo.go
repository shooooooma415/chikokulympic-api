package repository

import "chikokulympic-api/domain/entity"

type UserRepository interface {
	FindUserByUserID(userID entity.UserID) (*entity.User, error)
	FindUserByAuthID(authID entity.AuthID) (*entity.User, error)
	CreateUser(user entity.User) (*entity.User, error)
	DeleteUser(userID entity.UserID) (*entity.User, error)
	UpdateUser(user entity.User) (*entity.User, error)
}
