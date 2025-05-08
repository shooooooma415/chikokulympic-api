package entity

type Latitude float64
type Longitude float64

type UserLocation struct {
	UserID    UserID    `bson:"user_id" json:"user_id"`
	Latitude  Latitude  `bson:"latitude" json:"latitude"`
	Longitude Longitude `bson:"longitude" json:"longitude"`
}
