package location

import "chikokulympic-api/domain/user"

type LocationRepository interface {
	FindLocationByUserID(UserID user.UserID) (*UserLocation, error)
	CreateLocation(Location *UserLocation) (*UserLocation, error)
	DeleteLocation(Location *UserLocation) (*UserLocation, error)
	UpdateLocation(Location *UserLocation) (*UserLocation, error)
}
