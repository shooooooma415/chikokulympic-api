package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type CreateEventUseCase interface {
	Execute() (*entity.Event, error)
}

type CreateEventUseCaseImpl struct {
	eventRepo repository.EventRepository
	groupRepo repository.GroupRepository
	event     *entity.Event
	groupID   entity.GroupID
}

func NewCreateEventUseCase(eventRepo repository.EventRepository, groupRepo repository.GroupRepository, event *entity.Event, groupID entity.GroupID) *CreateEventUseCaseImpl {
	return &CreateEventUseCaseImpl{
		eventRepo: eventRepo,
		groupRepo: groupRepo,
		event:     event,
		groupID:   groupID,
	}
}

func (uc *CreateEventUseCaseImpl) Execute() (*entity.Event, error) {
	// イベントを作成
	createdEvent, err := uc.eventRepo.CreateEvent(*uc.event)
	if err != nil {
		return nil, err
	}

	// 指定したグループを取得
	group, err := uc.groupRepo.FindGroupByGroupID(uc.groupID)
	if err != nil {
		return nil, err
	}

	// グループのイベントリストにイベントIDを追加
	group.GroupEvents = append(group.GroupEvents, createdEvent.EventID)

	// グループを更新
	_, err = uc.groupRepo.UpdateGroup(*group)
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}
