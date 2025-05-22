package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
	"sort"
	"time"
)

type ArrivalRank struct {
	Rank        int             `json:"rank"`
	UserID      entity.UserID   `json:"user_id"`
	Name        entity.UserName `json:"name"`
	Alias       entity.Alias    `json:"alias"`
	ArrivalTime int             `json:"arrival_ime"`
}

type GetArrivalRankingResponse struct {
	EventID entity.EventID `json:"event_id"`
	Ranking []ArrivalRank  `json:"ranking"`
}

type GetArrivalRankingUseCase interface {
	Execute() (*GetArrivalRankingResponse, error)
}

type GetArrivalRankingUseCaseImpl struct {
	eventRepo repository.EventRepository
	eventID   *entity.EventID
	userRepo  repository.UserRepository
}

func NewGetArrivalRankingUseCase(eventRepo repository.EventRepository, eventID *entity.EventID, userRepo repository.UserRepository) *GetArrivalRankingUseCaseImpl {
	return &GetArrivalRankingUseCaseImpl{
		eventRepo: eventRepo,
		eventID:   eventID,
		userRepo:  userRepo,
	}
}

func (uc *GetArrivalRankingUseCaseImpl) Execute() (*GetArrivalRankingResponse, error) {
	event, err := uc.eventRepo.FindEventByEventID(*uc.eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	if len(event.VotedMembers) == 0 {
		return &GetArrivalRankingResponse{
			EventID: event.EventID,
			Ranking: []ArrivalRank{},
		}, nil
	}

	var userIDs []entity.UserID
	for _, member := range event.VotedMembers {
		if !member.IsArrival {
			continue // 到着していないメンバーはスキップ
		}
		userIDs = append(userIDs, member.UserID)
	}

	userMap := make(map[entity.UserID]*entity.User)
	for _, userID := range userIDs {
		user, err := uc.userRepo.FindUserByUserID(userID)
		if err != nil {
			continue
		}
		if user != nil {
			userMap[userID] = user
		}
	}

	eventStartTime := time.Time(event.EventStartDateTime)

	// ランキングの作成
	var ranking []ArrivalRank
	for _, member := range event.VotedMembers {
		if !member.IsArrival {
			continue // 参加していないメンバーはスキップ
		}

		// ユーザー情報を取得
		user, exists := userMap[member.UserID]
		if !exists {
			continue // ユーザー情報がない場合はスキップ
		}

		timeDiff := member.ArrivalDateTime.Sub(eventStartTime)
		timeDiffMinutes := int(timeDiff.Minutes())

		// ランキングエントリーを作成
		rankEntry := ArrivalRank{
			UserID:      member.UserID,
			Name:        user.UserName,
			Alias:       user.Alias,
			ArrivalTime: timeDiffMinutes,
		}
		ranking = append(ranking, rankEntry)
	}

	// ソート
	sort.Slice(ranking, func(i, j int) bool {
		return ranking[i].ArrivalTime < ranking[j].ArrivalTime
	})

	// ランク付け
	for i := range ranking {
		ranking[i].Rank = i + 1
	}

	return &GetArrivalRankingResponse{
		EventID: event.EventID,
		Ranking: ranking,
	}, nil
}
