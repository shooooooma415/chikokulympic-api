package group

import "chikokulympic-api/domain/user"

type GroupRepository interface {
	FindGroupByGroupID(GroupID GroupID) (*Group, error)
	FindGroupByUserID(UserID user.UserID) (*Group, error)
	CreateGroup(Group *Group) (*Group, error)
	DeleteGroup(Group *Group) (*Group, error)
	UpdateGroup(Group *Group) (*Group, error)
}
