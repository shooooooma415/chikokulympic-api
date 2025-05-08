package entity

type GroupID string
type GroupName string
type GroupPassword string
type GroupDescription string

type GroupMembers []UserID
type GroupEvents []EventID

type Group struct {
	GroupID          GroupID          `bson:"group_id" json:"group_id"`
	GroupName        GroupName        `bson:"name" json:"group_name"`
	GroupPassword    GroupPassword    `bson:"password" json:"group_password"`
	GroupManagerID   UserID           `bson:"manager_id" json:"group_manager_id"`
	GroupDescription GroupDescription `bson:"description" json:"group_description"`
	GroupMembers     GroupMembers     `bson:"members" json:"group_members"`
	GroupEvents      GroupEvents      `bson:"events" json:"group_events"`
}
