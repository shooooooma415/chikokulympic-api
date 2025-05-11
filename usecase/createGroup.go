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
	userRepo  repository.UserRepository
	group     *entity.Group
}

func NewCreateGroupUseCase(groupRepo repository.GroupRepository, userRepo repository.UserRepository, group *entity.Group) *CreateGroupUseCaseImpl {
	return &CreateGroupUseCaseImpl{
		groupRepo: groupRepo,
		userRepo:  userRepo,
		group:     group,
	}
}

func (uc *CreateGroupUseCaseImpl) Execute() (*entity.Group, error) {
	user, err := uc.userRepo.FindUserByUserID(uc.group.GroupManagerID)
	if err != nil {
		return nil, fmt.Errorf("ユーザー検索中にエラーが発生しました: %v", err)
	}
	if user == nil {
		return nil, fmt.Errorf("指定されたユーザーID %s が存在しません", string(uc.group.GroupManagerID))
	}

	existingGroup, err := uc.groupRepo.FindGroupByGroupName(&uc.group.GroupName)
	if err == nil && existingGroup != nil {
		return nil, fmt.Errorf("グループ名 '%s' は既に使用されています", string(uc.group.GroupName))
	}

	uc.group.GroupMembers = append(uc.group.GroupMembers, uc.group.GroupManagerID)

	return uc.groupRepo.CreateGroup(uc.group)
}
