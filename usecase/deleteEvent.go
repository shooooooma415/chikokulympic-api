package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

type DeleteEventUseCase interface {
	Execute() (*entity.Event, error)
}

type DeleteEventUseCaseImpl struct {
	eventRepo repository.EventRepository
	eventID   *entity.EventID
	authID    *entity.UserID
}

func NewDeleteEventUseCase(eventRepo repository.EventRepository, eventID *entity.EventID, authID *entity.UserID) *DeleteEventUseCaseImpl {
	return &DeleteEventUseCaseImpl{
		eventRepo: eventRepo,
		eventID:   eventID,
		authID:    authID,
	}
}
func (uc *DeleteEventUseCaseImpl) Execute() (*entity.Event, error) {
	event, err := uc.eventRepo.FindEventByEventID(*uc.eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	// イベントの作成者と認証ユーザーが一致するか確認
	if event.EventAuthorID != *uc.authID {
		return nil, fmt.Errorf("not authorized to delete this event")
	}

	return uc.eventRepo.DeleteEvent(*event)
}