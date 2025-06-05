package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/middleware"
	"chikokulympic-api/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PostEventRequest struct {
	GroupID              entity.GroupID              `json:"group_id" example:"group123"`
	EventID              entity.EventID              `json:"event_id" example:"event123"`
	EventTitle           entity.EventTitle           `json:"event_title" example:"テストイベント"`
	EventDescription     entity.EventDescription     `json:"event_description" example:"これはテストイベントです"`
	EventLocationName    entity.LocationName         `json:"event_location_name" example:"東京ドーム"`
	Cost                 entity.Cost                 `json:"cost" example:"1000"`
	EventMessage         entity.EventMessage         `json:"event_message" example:"参加してください！"`
	EventAuthorID        entity.UserID               `json:"event_author_id" example:"user123"`
	Latitude             entity.Latitude             `json:"latitude" example:"35.6895"`
	Longitude            entity.Longitude            `json:"longitude" example:"139.6917"`
	EventStartDateTime   entity.StartDateTIme        `json:"event_start_date_time" example:"2023-10-01T10:00:00Z"`
	EventEndDateTime     entity.EndDateTime          `json:"event_end_date_time" example:"2023-10-01T12:00:00Z"`
	EventClosingDateTime entity.EventClosingDateTime `json:"event_closing_date_time" example:"2023-09-30T23:59:59Z"`
}

type PostEventResponse struct {
	EventID entity.EventID `json:"event_id" example:"event123"`
}

type PostEvent struct {
	groupRepo repository.GroupRepository
	eventRepo repository.EventRepository
}

func NewPostEvent(groupRepo repository.GroupRepository, eventRepo repository.EventRepository) *PostEvent {
	return &PostEvent{
		groupRepo: groupRepo,
		eventRepo: eventRepo,
	}
}

// @Summary create event
// @Description create a new event
// @Tags events
// @Accept json
// @Produce json
// @Param request body PostEventRequest true "request"
// @Success 201 {object} PostEventResponse
// @Failure 400 {object} middleware.ErrorResponse
// @Failure 500 {object} middleware.ErrorResponse
// @Router /events [post]

func (p *PostEvent) Execute(c echo.Context, req PostEventRequest) error {
	event := &entity.Event{
		EventID:              req.EventID,
		EventTitle:           req.EventTitle,
		EventDescription:     req.EventDescription,
		EventLocationName:    req.EventLocationName,
		Cost:                 req.Cost,
		EventMessage:         req.EventMessage,
		EventAuthorID:        req.EventAuthorID,
		Latitude:             req.Latitude,
		Longitude:            req.Longitude,
		EventStartDateTime:   req.EventStartDateTime,
		EventEndDateTime:     req.EventEndDateTime,
		EventClosingDateTime: req.EventClosingDateTime,
	}

	createdEvent, err := usecase.NewCreateEventUseCase(p.eventRepo, p.groupRepo, event, req.GroupID).Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, middleware.NewErrorResponse(err.Error()))
	}

	response := &PostEventResponse{
		EventID: createdEvent.EventID,
	}

	return c.JSON(http.StatusCreated, response)
}
