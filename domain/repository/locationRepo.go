package repository

import "chikokulympic-api/domain/entity"

type LocationRepository interface {
	FindLocationByUserID(UserID entity.UserID) (*entity.UserLocation, error)
	CreateLocation(Location entity.UserLocation) (*entity.UserLocation, error)
	DeleteLocation(Location entity.UserLocation) (*entity.UserLocation, error)
	UpdateLocation(Location entity.UserLocation) (*entity.UserLocation, error)
}
