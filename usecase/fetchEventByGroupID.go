package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

type FetchEventInfoUsecase interface {
	Execute() (*entity.Event, error)
}
type FetchGroupInfoUsecaseImpl struct {
	eventRepo repository.EventRepository
	eventID   *entity.EventID
}

func NewFetchEventInfoUsecase(eventRepo repository.EventRepository, eventID *entity.EventID) *FetchGroupInfoUsecaseImpl {
	return &FetchGroupInfoUsecaseImpl{
		eventRepo: eventRepo,
		eventID:   eventID,
	}
}

func (uc *FetchGroupInfoUsecaseImpl) Execute() (*entity.Event, error) {
	event, err := uc.eventRepo.FindEventByEventID(*uc.eventID)
	if err != nil {
		return nil, fmt.Errorf("グループ情報の取得中にエラーが発生しました: %v", err)
	}
	if event == nil {
		return nil, fmt.Errorf("指定されたイベントID %s が存在しません", *uc.eventID)
	}
	return event, nil
}
