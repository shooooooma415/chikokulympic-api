package location

import user "chikokulympic-api/domain/auth"

type Latitude float64
type Longitude float64

type UserLocation struct {
	UserID    user.UserID
	Latitude  Latitude  `json:"latitude"`
	Longitude Longitude `json:"longitude"`
}
