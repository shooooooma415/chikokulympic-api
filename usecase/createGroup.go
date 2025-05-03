package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

type CreateGroupUseCase interface {
	Execute(group *entity.Group) (*entity.Group, error)
}

type CreateGroupUseCaseImpl struct {
	groupRepo repository.GroupRepository
}

func NewCreateGroupUseCase(groupRepo repository.GroupRepository) *CreateGroupUseCaseImpl {
	return &CreateGroupUseCaseImpl{
		groupRepo: groupRepo,
	}
}

func (uc *CreateGroupUseCaseImpl) Execute(group *entity.Group) (*entity.Group, error) {
	existingGroup, err := uc.groupRepo.FindGroupByGroupName(&group.GroupName)
	if err == nil && existingGroup != nil {
		return nil, fmt.Errorf("グループ名 '%s' は既に使用されています", string(group.GroupName))
	}

	return uc.groupRepo.CreateGroup(group)
}
