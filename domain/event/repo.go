package event

type EventRepository interface {
	FindEventByEventID(EventID EventID) (*Event, error)
	CreateEvent(Event *Event) (*Event, error)
	DeleteEvent(Event *Event) (*Event, error)
	UpdateEvent(Event *Event) (*Event, error)
}
