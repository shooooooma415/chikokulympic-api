package entity

import (
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
	EventAuthorID        UserID
	Latitude             Latitude
	Longitude            Longitude
	EventStartDateTime   StartDateTIme
	EventEndDateTime     EndDateTime
	EventClosingDateTime EventClosingDateTime
}
