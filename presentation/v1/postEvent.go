package v1

import "chikokulympic-api/domain/entity"

type PostEventRequest struct {
	GroupID              entity.GroupID          `json:"group_id" example:"group123"`
	EventID              entity.EventID          `json:"event_id" example:"event123"`
	EventTitle           entity.EventTitle       `json:"event_title" example:"テストイベント"`
	EventDescription     entity.EventDescription `json:"event_description" example:"これはテストイベントです"`
	EventLocationName    entity.LocationName     `json:"event_location_name" example:"東京ドーム"`
	Cost                 entity.Cost             `json:"cost" example:"1000"`
	EventMessage         entity.EventMessage     `json:"event_message" example:"参加してください！"`
	EventAuthorID        entity.UserID           `json:"event_author_id" example:"user123"`
	Latitude             entity.Latitude         `json:"latitude" example:"35.6895"`
	Longitude            entity.Longitude        `json:"longitude" example:"139.6917"`
	EventStartDateTime   string                  `json:"event_start_date_time" example:"2023-10-01T10:00:00Z"`
	EventEndDateTime     string                  `json:"event_end_date_time" example:"2023-10-01T12:00:00Z"`
	EventClosingDateTime string                  `json:"event_closing_date_time" example:"2023-09-30T23:59:59Z"`
}

type PostEventResponse struct {
	EventID              entity.EventID          `json:"event_id" example:"event123"`
}