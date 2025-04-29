package user

type UserRepository interface {
	FindUserByUserID(UserID UserID) (*User, error)
	FindUserByAuthID(AuthID AuthID) (*User, error)
	CreateUser(User *User) (*User, error)
	DeleteUser(UserID UserID) (*User, error)
	UpdateUser(User *User) (*User, error)
}
