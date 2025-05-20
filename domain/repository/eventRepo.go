package repository

import "chikokulympic-api/domain/entity"

type EventRepository interface {
	FindEventByEventID(eventID entity.EventID) (*entity.Event, error)
	CreateEvent(event entity.Event) (*entity.Event, error)
	DeleteEvent(event entity.Event) (*entity.Event, error)
	UpdateEvent(event entity.Event) (*entity.Event, error)
}
