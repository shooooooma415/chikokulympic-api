package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
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
	events := make([]entity.Event, 0)

	for _, groupID := range uc.groupIDs {
		group, err := uc.groupRepo.FindGroupByGroupID(groupID)
		if err != nil {
			return nil, fmt.Errorf("グループが見つかりません: %v", err)
		}
		for _, eventID := range group.GroupEvents {
			event, err := uc.eventRepo.FindEventByEventID(eventID)
			if err != nil {
				return nil, fmt.Errorf("イベントが見つかりません: %v", err)
			}
			if event == nil {
				return nil, fmt.Errorf("イベントが見つかりません")
			}
			events = append(events, *event)
		}

	}
	return events, nil
}
