package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type CreateGroupUseCase interface {
	Execute(groupName string, groupPassword string, groupManagerID string, groupDescription string) (*entity.Group, error)
}

type CreateGroupUseCaseImpl struct {
	groupRepo repository.GroupRepository
}

func NewCreateGroupUseCase(groupRepo repository.GroupRepository) *CreateGroupUseCaseImpl {
	return &CreateGroupUseCaseImpl{
		groupRepo: groupRepo,
	}
}
func (uc *CreateGroupUseCaseImpl) Execute(groupName string, groupPassword string, groupManagerID string, groupDescription string) (*entity.Group, error) {
	group := &entity.Group{
		GroupName:        entity.GroupName(groupName),
		GroupPassword:    entity.GroupPassword(groupPassword),
		GroupManagerID:   entity.UserID(groupManagerID),
		GroupDescription: entity.GroupDescription(groupDescription),
	}

	return uc.groupRepo.CreateGroup(group)
}
