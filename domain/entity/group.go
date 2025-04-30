package entity

type GroupID string
type GroupName string
type GroupDescription string

type GroupMembers []UserID
type GroupEvents []EventID

type Group struct {
	GroupID          GroupID          `bson:"group_id" json:"group_id"`
	GroupName        GroupName        `bson:"group_name" json:"group_name"`
	GroupManagerID   UserID           `bson:"group_manager_id" json:"group_manager_id"`
	GroupDescription GroupDescription `bson:"group_description" json:"group_description"`
	GroupMembers     GroupMembers     `bson:"group_members" json:"group_members"`
	GroupEvents      GroupEvents      `bson:"group_events" json:"group_events"`
}
