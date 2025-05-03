package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

type JoinGroupUseCase interface {
	Execute(userID entity.UserID, group entity.Group) (*entity.GroupID, error)
}

type JoinGroupUseCaseImpl struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

func NewJoinGroupUseCase(groupRepo repository.GroupRepository, userRepo repository.UserRepository) *JoinGroupUseCaseImpl {
	return &JoinGroupUseCaseImpl{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

func (uc *JoinGroupUseCaseImpl) Execute(userID entity.UserID, group entity.Group) (*entity.GroupID, error) {
	groupFound, err := uc.groupRepo.FindGroupByGroupName(&group.GroupName)
	if err != nil {
		return nil, err
	}
	if groupFound == nil {
		return nil, fmt.Errorf("グループが見つかりません")
	}

	if groupFound.GroupPassword != group.GroupPassword {
		return nil, fmt.Errorf("パスワードが一致しません")
	}

	user, err := uc.userRepo.FindUserByUserID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("ユーザーが見つかりません")
	}

	for _, memberID := range groupFound.GroupMembers {
		if memberID == user.UserID {
			return nil, fmt.Errorf("すでにグループに参加しています")
		}
	}

	if groupFound.GroupManagerID == user.UserID {
		return nil, fmt.Errorf("あなたはこのグループのマネージャーです")
	}

	groupFound.GroupMembers = append(groupFound.GroupMembers, user.UserID)

	updatedGroup, err := uc.groupRepo.UpdateGroup(groupFound)
	if err != nil {
		return nil, err
	}

	return &updatedGroup.GroupID, nil
}
