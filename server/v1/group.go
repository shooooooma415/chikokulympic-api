package v1

import (
	presentationV1 "chikokulympic-api/presentation/v1"
	"chikokulympic-api/domain/repository"
	"github.com/labstack/echo/v4"
)

type GroupServer struct {
	joinGroup     *presentationV1.JoinGroup
	leaveGroup    *presentationV1.LeaveGroup
	getUserGroups *presentationV1.GetUserGroups
	getGroupInfo  *presentationV1.GetGroupInfo
}

func NewGroupServer(groupRepo repository.GroupRepository, userRepo repository.UserRepository) *GroupServer {
	return &GroupServer{
		joinGroup:     presentationV1.NewJoinGroup(userRepo, groupRepo),
		leaveGroup:    presentationV1.NewLeaveGroup(groupRepo),
		getUserGroups: presentationV1.NewGetUserGroups(groupRepo),
		getGroupInfo:  presentationV1.NewGetGroupInfo(groupRepo, userRepo),
	}
}
func (s *GroupServer) RegisterRoutes(e *echo.Echo) {
	e.POST("/group/join", s.joinGroup.Handler)

	e.POST("/group/leave", s.leaveGroup.Handler)

	e.GET("/user/:user_id/groups", s.getUserGroups.Handler)

	e.GET("/group/:group_id", s.getGroupInfo.Handler)
}
