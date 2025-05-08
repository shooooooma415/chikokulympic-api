package repository

import (
	"chikokulympic-api/domain/entity"
)

type GroupRepository interface {
	FindGroupByGroupName(Group *entity.GroupName) (*entity.Group, error)
	FindGroupsByUserID(UserID entity.UserID) ([]*entity.Group, error)
	CreateGroup(Group *entity.Group) (*entity.Group, error)
	DeleteGroup(Group *entity.Group) (*entity.Group, error)
	UpdateGroup(Group *entity.Group) (*entity.Group, error)
}
