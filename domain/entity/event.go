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
type Vote string

type VotedMember struct {
	IsArrival       bool      `bson:"is_arrival"`
	UserID          UserID    `bson:"user_id"`
	ArrivalDateTime time.Time `bson:"arrival_date_time"`
}

type Event struct {
	EventID              EventID              `bson:"event_id"`
	EventTitle           EventTitle           `bson:"event_title"`
	EventDescription     EventDescription     `bson:"event_description"`
	EventLocationName    LocationName         `bson:"event_location_name"`
	Cost                 Cost                 `bson:"cost"`
	EventMessage         EventMessage         `bson:"event_message"`
	EventAuthorID        UserID               `bson:"event_author_id"`
	Latitude             Latitude             `bson:"latitude"`
	Longitude            Longitude            `bson:"longitude"`
	EventStartDateTime   StartDateTIme        `bson:"event_start_date_time"`
	EventEndDateTime     EndDateTime          `bson:"event_end_date_time"`
	EventClosingDateTime EventClosingDateTime `bson:"event_closing_date_time"`
	VotedMembers         []VotedMember        `bson:"voted_members"`
}
