package repository

import "chikokulympic-api/domain/entity"

type LocationRepository interface {
	FindLocationByUserID(userID entity.UserID) (*entity.UserLocation, error)
	CreateLocation(location entity.UserLocation) (*entity.UserLocation, error)
	DeleteLocation(location entity.UserLocation) (*entity.UserLocation, error)
	UpdateLocation(location entity.UserLocation) (*entity.UserLocation, error)
}
