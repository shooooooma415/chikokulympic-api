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
	createdEvent, err := uc.eventRepo.CreateEvent(*uc.event)
	if err != nil {
		return nil, err
	}

	group, err := uc.groupRepo.FindGroupByGroupID(uc.groupID)
	if err != nil {
		return nil, err
	}

	group.GroupEvents = append(group.GroupEvents, createdEvent.EventID)

	_, err = uc.groupRepo.UpdateGroup(*group)
	if err != nil {
		return nil, err
	}

	return createdEvent, nil
}
