package group

type GroupRepository interface {
	FindGroupByGroupID(GroupID GroupID) (*Group, error)
	CreateGroup(Group *Group) (*Group, error)
	DeleteGroup(Group *Group) (*Group, error)
	UpdateGroup(Group *Group) (*Group, error)
}
