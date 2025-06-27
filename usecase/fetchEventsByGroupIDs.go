package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
	"sort"
	"sync"
	"time"
)

type FetchEventsByGroupIDsUsecase interface {
	Execute() ([]entity.Event, error)
}
type FetchEventsByGroupIDsUsecaseImpl struct {
	groupRepo repository.GroupRepository
	eventRepo repository.EventRepository
	groupIDs  []entity.GroupID
}

func NewFetchEventInfoUsecase(groupRepo repository.GroupRepository, eventRepo repository.EventRepository, groupIDs []entity.GroupID) *FetchEventsByGroupIDsUsecaseImpl {
	return &FetchEventsByGroupIDsUsecaseImpl{
		groupRepo: groupRepo,
		eventRepo: eventRepo,
		groupIDs:  groupIDs,
	}
}

func (uc *FetchEventsByGroupIDsUsecaseImpl) Execute() ([]entity.Event, error) {
	var (
		events     []entity.Event
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

				mutex.Lock()
				events = append(events, *event)
				mutex.Unlock()
			}
		}(groupID)
	}

	wg.Wait()

	if errOccured {
		return nil, firstErr
	}

	// EventStartDateTimeの降順でソート
	sort.Slice(events, func(i, j int) bool {
		// time.Timeに変換して比較（降順なのでj > iの順番で比較）
		timeI := time.Time(events[i].EventStartDateTime)
		timeJ := time.Time(events[j].EventStartDateTime)
		return timeJ.Before(timeI)
	})

	return events, nil
}
