package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

type CreateGroupUseCase interface {
	Execute() (*entity.Group, error)
}

type CreateGroupUseCaseImpl struct {
	groupRepo repository.GroupRepository
	group     *entity.Group
}

func NewCreateGroupUseCase(groupRepo repository.GroupRepository, group *entity.Group) *CreateGroupUseCaseImpl {
	return &CreateGroupUseCaseImpl{
		groupRepo: groupRepo,
		group:     group,
	}
}

func (uc *CreateGroupUseCaseImpl) Execute() (*entity.Group, error) {
	existingGroup, err := uc.groupRepo.FindGroupByGroupName(&uc.group.GroupName)
	if err == nil && existingGroup != nil {
		return nil, fmt.Errorf("グループ名 '%s' は既に使用されています", string(uc.group.GroupName))
	}

	return uc.groupRepo.CreateGroup(uc.group)
}
