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

// レスポンス構造体
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
	// イベントを取得
	event, err := uc.eventRepo.FindEventByEventID(*uc.eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("イベントが見つかりません")
	}

	// イベントに投票済みのメンバー情報を処理
	if len(event.VotedMembers) == 0 {
		return &GetArrivalRankingResponse{
			EventID: event.EventID,
			Ranking: []ArrivalRank{},
		}, nil
	}

	// ユーザーIDのリストを取得して、ユーザー情報を取得
	var userIDs []entity.UserID
	for _, member := range event.VotedMembers {
		if !member.IsArrival {
			continue // 参加していないメンバーはスキップ
		}
		userIDs = append(userIDs, member.UserID)
	}

	// ユーザー情報をマップに格納
	userMap := make(map[entity.UserID]*entity.User)
	for _, userID := range userIDs {
		// 各ユーザー情報を個別に取得
		user, err := uc.userRepo.FindUserByUserID(userID)
		if err != nil {
			// エラーが発生してもスキップして続行
			continue
		}
		if user != nil {
			userMap[userID] = user
		}
	}

	// イベント開始時間
	eventStartTime := time.Time(event.EventStartDateTime)

	// 到着時間とイベント開始時間の差でランキングを作成
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

		// 時間差（ミリ秒）を計算
		timeDiff := member.ArrivalDateTime.Sub(eventStartTime)

		// 分単位に変換
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

	// 時間差が小さい順（差が少ないほど高ランク）にソート
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
