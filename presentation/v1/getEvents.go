package v1

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"chikokulympic-api/usecase"
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
	"chikokulympic-api/middleware"
)
type GetEventsResponse struct {
	Events []entity.Event `json:"events" example:"[{\"event_id\":\"event123\",\"event_name\":\"Event Name\",\"event_description\":\"Event Description\",\"event_date\":\"2023-10-01T00:00:00Z\"}]"`
}

type GetEvents struct {
	eventRepo repository.EventRepository
	groupRepo repository.GroupRepository
}

func NewGetEvents(eventRepo repository.EventRepository, groupRepo repository.GroupRepository) *GetEvents {
	return &GetEvents{
		eventRepo: eventRepo,
		groupRepo: groupRepo,
	}
}

