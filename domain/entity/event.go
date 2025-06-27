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
	IsArrival       bool      `bson:"is_arrival" json:"is_arrival"`
	UserID          UserID    `bson:"user_id" json:"user_id"`
	Vote            Vote      `bson:"vote" json:"vote"`
	ArrivalDateTime time.Time `bson:"arrival_date_time" json:"arrival_date_time"`
}

type Event struct {
	EventID              EventID              `bson:"_id" json:"event_id"`
	EventTitle           EventTitle           `bson:"event_title" json:"event_title"`
	EventDescription     EventDescription     `bson:"event_description" json:"event_description"`
	EventLocationName    LocationName         `bson:"event_location_name" json:"event_location_name"`
	Cost                 Cost                 `bson:"cost" json:"cost"`
	EventMessage         EventMessage         `bson:"event_message" json:"event_message"`
	EventAuthorID        UserID               `bson:"event_author_id" json:"event_author_id"`
	Latitude             Latitude             `bson:"latitude" json:"latitude"`
	Longitude            Longitude            `bson:"longitude" json:"longitude"`
	EventStartDateTime   StartDateTIme        `bson:"event_start_date_time" json:"event_start_date_time"`
	EventEndDateTime     EndDateTime          `bson:"event_end_date_time" json:"event_end_date_time"`
	EventClosingDateTime EventClosingDateTime `bson:"event_closing_date_time" json:"event_closing_date_time"`
	VotedMembers         []VotedMember        `bson:"voted_members" json:"voted_members"`
}
