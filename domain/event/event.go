package event

import (
	"chikokulympic-api/domain/location"
	"chikokulympic-api/domain/user"
	"time"
)

type EventID string
type EventTitle string
type EventDescription string
type LocationName string
type Cost int
type EventMessage string
type StartDateTIme time.Time
type EndDateTime time.Time
type EventClosingDateTime time.Time

type Event struct {
	EventID              EventID
	EventTitle           EventTitle
	EventDescription     EventDescription
	EventLocationName    LocationName
	Cost                 Cost
	EventMessage         EventMessage
	EventAuthorID        user.UserID
	Latitude             location.Latitude
	Longitude            location.Longitude
	EventStartDateTime   StartDateTIme
	EventEndDateTime     EndDateTime
	EventClosingDateTime EventClosingDateTime
}
