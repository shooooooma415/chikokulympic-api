package repository

import "chikokulympic-api/domain/entity"

type EventRepository interface {
	FindEventByEventID(EventID entity.EventID) (*entity.Event, error)
	CreateEvent(Event *entity.Event) (*entity.Event, error)
	DeleteEvent(Event *entity.Event) (*entity.Event, error)
	UpdateEvent(Event *entity.Event) (*entity.Event, error)
}
