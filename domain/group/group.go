package group

import (
	"chikokulympic-api/domain/user"
)

type GroupID string
type GroupName string
type GroupDescription string

type GroupMembers []user.UserID

type Group struct {
	GroupID          GroupID
	GroupName        GroupName
	GroupDescription GroupDescription
	GroupMembers     GroupMembers
}
