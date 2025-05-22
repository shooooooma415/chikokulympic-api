package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"time"
)

type ArrivalRank struct {
	Rank        int           `json:"rank"`
	UserID      entity.UserID `json:"user_id"`
	Name        string        `json:"name"`
	Alias       entity.Alias  `json:"alias"`
	ArrivalTime time.Time     `json:"arrival_time"`
}

type getArrivalRankingResponse struct {
	EventID entity.EventID
	Ranking []ArrivalRank
}3

type getArrivalRankingUseCase interface {
	Execute() ([]getArrivalRankingResponse, error)
}

type getArrivalRankingUseCaseImpl struct {
	eventRepo repository.EventRepository
	eventID   *entity.EventID
	userRepo  repository.UserRepository
}

func NewGetArrivalRankingUseCase(eventRepo repository.EventRepository, eventID *entity.EventID, userRepo repository.UserRepository) *getArrivalRankingUseCaseImpl {
	return &getArrivalRankingUseCaseImpl{
		eventRepo: eventRepo,
		eventID:   eventID,
		userRepo:  userRepo,
	}
}
