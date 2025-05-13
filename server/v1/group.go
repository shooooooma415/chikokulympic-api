package v1

import (
	presentationV1 "chikokulympic-api/presentation/v1"
	"chikokulympic-api/domain/repository"
	"github.com/labstack/echo/v4"
)

type GroupServer struct {
	createGroup   *presentationV1.PostGroup
	joinGroup     *presentationV1.JoinGroup
	leaveGroup    *presentationV1.LeaveGroup
	getGroupInfo  *presentationV1.GetGroupInfo
}

func NewGroupServer(groupRepo repository.GroupRepository, userRepo repository.UserRepository) *GroupServer {
	return &GroupServer{
		createGroup:  presentationV1.NewPostGroup(groupRepo, userRepo),
		joinGroup:     presentationV1.NewJoinGroup(userRepo, groupRepo),
		leaveGroup:    presentationV1.NewLeaveGroup(groupRepo),
		getGroupInfo:  presentationV1.NewGetGroupInfo(groupRepo, userRepo),
	}
}
func (s *GroupServer) RegisterRoutes(e *echo.Echo) {
	groupGroup := e.Group("/groups")

	groupGroup.POST("", s.createGroup.Handler)

	groupGroup.POST("/join", s.joinGroup.Handler)

	groupGroup.POST("/:group_id/leave", s.leaveGroup.Handler)


	groupGroup.GET("/:group_id", s.getGroupInfo.Handler)
}
