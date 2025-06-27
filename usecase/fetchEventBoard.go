package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
	"sync"
	"time"
)

type EventBoardAuthor struct {
	AuthorID   string `json:"author_id"`
	AuthorName string `json:"author_name"`
}

type EventBoardOption struct {
	Title            string `json:"title"`
	ParticipantCount int    `json:"participant_count"`
	Participants     []struct {
		UserID   string `json:"user_id"`
		UserName string `json:"user_name"`
	} `json:"participants"`
}

type EventBoardEvent struct {
	ID           string             `json:"id"`
	Title        string             `json:"title"`
	Author       EventBoardAuthor   `json:"author"`
	Description  string             `json:"description"`
	IsAllDay     bool               `json:"is_all_day"`
	StartTime    time.Time          `json:"start_time"`
	EndTime      time.Time          `json:"end_time"`
	ClosingTime  time.Time          `json:"closing_time"`
	LocationName string             `json:"location_name"`
	Cost         int                `json:"cost"`
	Message      string             `json:"message"`
	Latitude     string             `json:"latitude"`
	Longitude    string             `json:"longitude"`
	GroupID      string             `json:"group_id"`
	GroupName    string             `json:"group_name"`
	Options      []EventBoardOption `json:"options"`
}

type FetchEventBoardResponse struct {
	Events []EventBoardEvent `json:"events"`
}

type FetchEventBoardUseCase interface {
	Execute() (*FetchEventBoardResponse, error)
}

type FetchEventBoardUseCaseImpl struct {
	groupRepo repository.GroupRepository
	eventRepo repository.EventRepository
	userRepo  repository.UserRepository
	groupIDs  []entity.GroupID
}

func NewFetchEventBoardUseCase(groupRepo repository.GroupRepository, eventRepo repository.EventRepository, userRepo repository.UserRepository, groupIDs []entity.GroupID) *FetchEventBoardUseCaseImpl {
	return &FetchEventBoardUseCaseImpl{
		groupRepo: groupRepo,
		eventRepo: eventRepo,
		userRepo:  userRepo,
		groupIDs:  groupIDs,
	}
}

func (uc *FetchEventBoardUseCaseImpl) Execute() (*FetchEventBoardResponse, error) {
	var (
		events     []EventBoardEvent
		mutex      sync.Mutex
		wg         sync.WaitGroup
		errOccured bool
		firstErr   error
		errMutex   sync.Mutex
	)

	// 各グループIDを並列処理
	for _, groupID := range uc.groupIDs {
		wg.Add(1)
		go func(gID entity.GroupID) {
			defer wg.Done()

			// グループ情報を取得
			group, err := uc.groupRepo.FindGroupByGroupID(gID)
			if err != nil {
				errMutex.Lock()
				if !errOccured {
					errOccured = true
					firstErr = fmt.Errorf("グループが見つかりません: %v", err)
				}
				errMutex.Unlock()
				return
			}

			// グループ内の各イベントIDを処理
			for _, eventID := range group.GroupEvents {
				event, err := uc.eventRepo.FindEventByEventID(eventID)
				if err != nil {
					errMutex.Lock()
					if !errOccured {
						errOccured = true
						firstErr = fmt.Errorf("イベントが見つかりません: %v", err)
					}
					errMutex.Unlock()
					return
				}

				if event == nil {
					errMutex.Lock()
					if !errOccured {
						errOccured = true
						firstErr = fmt.Errorf("イベントが見つかりません")
					}
					errMutex.Unlock()
					return
				}

				// イベント作成者の情報を取得
				author, err := uc.userRepo.FindUserByUserID(event.EventAuthorID)
				if err != nil {
					errMutex.Lock()
					if !errOccured {
						errOccured = true
						firstErr = fmt.Errorf("イベント作成者の情報取得に失敗しました: %v", err)
					}
					errMutex.Unlock()
					return
				}

				// 投票オプションを処理
				var options []EventBoardOption
				voteCounts := make(map[string]int)
				voteParticipants := make(map[string][]struct {
					UserID   string `json:"user_id"`
					UserName string `json:"user_name"`
				})

				// 投票メンバーからオプションを抽出
				for _, member := range event.VotedMembers {
					vote := string(member.Vote)
					voteCounts[vote]++

					// 参加者情報を取得
					participant, err := uc.userRepo.FindUserByUserID(member.UserID)
					if err == nil && participant != nil {
						voteParticipants[vote] = append(voteParticipants[vote], struct {
							UserID   string `json:"user_id"`
							UserName string `json:"user_name"`
						}{
							UserID:   string(participant.UserID),
							UserName: string(participant.UserName),
						})
					}
				}

				// オプションを作成
				for vote, count := range voteCounts {
					option := EventBoardOption{
						Title:            vote,
						ParticipantCount: count,
						Participants:     voteParticipants[vote],
					}
					options = append(options, option)
				}

				// イベントボードイベントを作成
				boardEvent := EventBoardEvent{
					ID:           string(event.EventID),
					Title:        string(event.EventTitle),
					Description:  string(event.EventDescription),
					IsAllDay:     false, // 現在のエンティティにはこのフィールドがないため、デフォルト値を設定
					StartTime:    time.Time(event.EventStartDateTime),
					EndTime:      time.Time(event.EventEndDateTime),
					ClosingTime:  time.Time(event.EventClosingDateTime),
					LocationName: string(event.EventLocationName),
					Cost:         int(event.Cost),
					Message:      string(event.EventMessage),
					Latitude:     fmt.Sprintf("%f", event.Latitude),
					Longitude:    fmt.Sprintf("%f", event.Longitude),
					GroupID:      string(group.GroupID),
					GroupName:    string(group.GroupName),
					Options:      options,
				}

				if author != nil {
					boardEvent.Author = EventBoardAuthor{
						AuthorID:   string(author.UserID),
						AuthorName: string(author.UserName),
					}
				}

				mutex.Lock()
				events = append(events, boardEvent)
				mutex.Unlock()
			}
		}(groupID)
	}

	wg.Wait()

	if errOccured {
		return nil, firstErr
	}

	return &FetchEventBoardResponse{
		Events: events,
	}, nil
}
