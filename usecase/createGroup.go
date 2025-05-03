package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
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
	// グループ名の重複チェック
	groupNameObj := entity.GroupName(groupName)
	existingGroup, err := uc.groupRepo.FindGroupByGroupName(&groupNameObj)
	if err == nil && existingGroup != nil {
		return nil, fmt.Errorf("グループ名 '%s' は既に使用されています", groupName)
	}

	group := &entity.Group{
		GroupName:        entity.GroupName(groupName),
		GroupPassword:    entity.GroupPassword(groupPassword),
		GroupManagerID:   entity.UserID(groupManagerID),
		GroupDescription: entity.GroupDescription(groupDescription),
	}

	return uc.groupRepo.CreateGroup(group)
}
