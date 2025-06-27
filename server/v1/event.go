package v1

import (
	"chikokulympic-api/domain/repository"
	presentationV1 "chikokulympic-api/presentation/v1"

	"github.com/labstack/echo/v4"
)

type EventServer struct {
	postEvent     *presentationV1.PostEvent
	getEvents     *presentationV1.GetEvents
	getEventBoard *presentationV1.GetEventBoard
	postVote      *presentationV1.PostVote
}

func NewEventServer(eventRepo repository.EventRepository, groupRepo repository.GroupRepository, userRepo repository.UserRepository) *EventServer {
	return &EventServer{
		postEvent:     presentationV1.NewPostEvent(groupRepo, eventRepo),
		getEvents:     presentationV1.NewGetEvents(eventRepo, groupRepo),
		getEventBoard: presentationV1.NewGetEventBoard(groupRepo, eventRepo, userRepo),
		postVote:      presentationV1.NewPostVote(eventRepo, groupRepo, userRepo),
	}
}

func (s *EventServer) RegisterRoutes(e *echo.Echo) {
	eventGroup := e.Group("/events")

	eventGroup.POST("", s.postEvent.Handler)
	eventGroup.GET("", s.getEvents.Handler)
	eventGroup.GET("/board", s.getEventBoard.Handler)
	eventGroup.POST("/:event_id/votes", s.postVote.Handler)
}
