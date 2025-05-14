package entity

type GroupID string

type GroupName string

type GroupPassword string

type GroupDescription string

type GroupMembers []UserID

type GroupEvents []EventID

type Group struct {
	GroupID          GroupID          `bson:"group_id" json:"group_id" example:"group123"`
	GroupName        GroupName        `bson:"name" json:"group_name" example:"テストグループ"`
	GroupPassword    GroupPassword    `bson:"password" json:"group_password" example:"password123"`
	GroupManagerID   UserID           `bson:"manager_id" json:"group_manager_id" example:"user456"`
	GroupDescription GroupDescription `bson:"description" json:"group_description" example:"これはテストグループです"`
	GroupMembers     GroupMembers     `bson:"members" json:"group_members" example:"[\"user123\",\"user456\"]"`
	GroupEvents      GroupEvents      `bson:"events" json:"group_events" example:"[\"event123\",\"event456\"]"`
}
