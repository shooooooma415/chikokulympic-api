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
	event     *entity.Event
}

func NewCreateEventUseCase(eventRepo repository.EventRepository, event *entity.Event) *CreateEventUseCaseImpl {
	return &CreateEventUseCaseImpl{
		eventRepo: eventRepo,
		event:     event,
	}
}

func (uc *CreateEventUseCaseImpl) Execute() (*entity.Event, error) {
	return uc.eventRepo.CreateEvent(*uc.event)
}
