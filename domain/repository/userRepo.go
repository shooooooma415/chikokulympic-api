package repository

import "chikokulympic-api/domain/entity"

type UserRepository interface {
	FindUserByUserID(UserID entity.UserID) (*entity.User, error)
	FindUserByAuthID(AuthID entity.AuthID) (*entity.User, error)
	CreateUser(User *entity.User) (*entity.User, error)
	DeleteUser(UserID *entity.UserID) (*entity.User, error)
	UpdateUser(User *entity.User) (*entity.User, error)
}
