package entity

type GroupID string
type GroupName string
type GroupDescription string

type GroupMembers []UserID
type GroupEvents []EventID

type Group struct {
	GroupID          GroupID
	GroupName        GroupName
	GroupDescription GroupDescription
	GroupMembers     GroupMembers
	GroupEvents      GroupEvents
}
