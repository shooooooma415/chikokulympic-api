package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

// JoinGroupUseCase はグループ参加のユースケースインターフェース
// @Summary グループに参加するユースケース
// @Description ユーザーが指定されたグループに参加するためのユースケース
// @Tags group
type JoinGroupUseCase interface {
	Execute() (*entity.GroupID, error)
}

// JoinGroupUseCaseImpl はグループ参加ユースケースの実装
type JoinGroupUseCaseImpl struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
	userID    entity.UserID
	group     entity.Group
}

// NewJoinGroupUseCase は新しいJoinGroupUseCaseを作成する
func NewJoinGroupUseCase(groupRepo repository.GroupRepository, userRepo repository.UserRepository, userID entity.UserID, group entity.Group) *JoinGroupUseCaseImpl {
	return &JoinGroupUseCaseImpl{
		groupRepo: groupRepo,
		userRepo:  userRepo,
		userID:    userID,
		group:     group,
	}
}

// Execute はグループ参加処理を実行する
// @Summary グループ参加処理を実行する
// @Description ユーザーがグループに参加する処理を実行します
// @Success 200 {object} entity.GroupID "成功時はグループIDを返す"
// @Failure 400 {string} string "グループが存在しない場合や、パスワードが一致しない場合など"
// @Failure 404 {string} string "ユーザーが見つからない場合"
// @Failure 409 {string} string "すでにグループに参加している場合"
func (uc *JoinGroupUseCaseImpl) Execute() (*entity.GroupID, error) {
	groupFound, err := uc.groupRepo.FindGroupByGroupName(&uc.group.GroupName)
	if err != nil {
		return nil, err
	}
	if groupFound == nil {
		return nil, fmt.Errorf("グループが見つかりません")
	}

	if groupFound.GroupPassword != uc.group.GroupPassword {
		return nil, fmt.Errorf("パスワードが一致しません")
	}

	user, err := uc.userRepo.FindUserByUserID(uc.userID)
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
