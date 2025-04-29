package entity

type Latitude float64
type Longitude float64

type UserLocation struct {
	UserID    UserID
	Latitude  Latitude  `json:"latitude"`
	Longitude Longitude `json:"longitude"`
}
