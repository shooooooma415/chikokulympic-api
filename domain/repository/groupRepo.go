package repository

import "chikokulympic-api/domain/entity"

type GroupRepository interface {
	FindGroupByGroupName(groupName entity.GroupName) (*entity.Group, error)
	FindGroupByGroupID(groupID entity.GroupID) (*entity.Group, error)
	FindGroupsByUserID(userID entity.UserID) ([]*entity.Group, error)
	CreateGroup(group entity.Group) (*entity.Group, error)
	DeleteGroup(group entity.Group) (*entity.Group, error)
	UpdateGroup(group entity.Group) (*entity.Group, error)
}
