package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

// GroupResponse はグループ情報レスポンスの構造体
// @Description グループ情報レスポンス
type GroupResponse struct {
	ID          string `json:"id" example:"group123"`
	Name        string `json:"name" example:"テストグループ"`
	MemberCount int    `json:"member_count" example:"5"`
	IsCreator   bool   `json:"is_creator" example:"true"`
}

// FetchUserGroupsResponse はユーザーが所属するグループ一覧レスポンスの構造体
// @Description ユーザーが所属するグループ一覧レスポンス
type FetchUserGroupsResponse struct {
	Groups []GroupResponse `json:"groups"`
}

// UserGroup はユーザーグループのレスポンス型（APIドキュメント用）
// @Description ユーザーグループ情報
type UserGroup FetchUserGroupsResponse

// FetchUserGroupsUseCase はユーザーのグループ一覧取得インターフェース
type FetchUserGroupsUseCase interface {
	Execute() (*FetchUserGroupsResponse, error)
}

// FetchUserGroupsUseCaseImpl はユーザーのグループ一覧取得ユースケースの実装
type FetchUserGroupsUseCaseImpl struct {
	groupRepo repository.GroupRepository
	userID    entity.UserID
}

// NewFetchUserGroupsUseCase は新しいFetchUserGroupsUseCaseを作成する
func NewFetchUserGroupsUseCase(groupRepo repository.GroupRepository, userID entity.UserID) *FetchUserGroupsUseCaseImpl {
	return &FetchUserGroupsUseCaseImpl{
		groupRepo: groupRepo,
		userID:    userID,
	}
}

// Execute はユーザーのグループ一覧取得処理を実行する
func (uc *FetchUserGroupsUseCaseImpl) Execute() (*FetchUserGroupsResponse, error) {
	groups, err := uc.groupRepo.FindGroupsByUserID(uc.userID)
	if err != nil {
		return nil, err
	}

	result := &FetchUserGroupsResponse{
		Groups: make([]GroupResponse, 0, len(groups)),
	}

	for _, group := range groups {
		memberCount := len(group.GroupMembers)
		if group.GroupManagerID != "" {
			memberCount++
		}

		isCreator := group.GroupManagerID == uc.userID

		groupResponse := GroupResponse{
			ID:          string(group.GroupID),
			Name:        string(group.GroupName),
			MemberCount: memberCount,
			IsCreator:   isCreator,
		}

		result.Groups = append(result.Groups, groupResponse)
	}

	return result, nil
}
