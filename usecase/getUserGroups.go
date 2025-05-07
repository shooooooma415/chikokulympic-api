package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type GroupInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"member_count"`
	IsCreator   bool   `json:"is_creator"`
}

type GetUserGroupsResponse struct {
	Groups []GroupInfo `json:"groups"`
}

type GetUserGroupsUseCase interface {
	Execute() (*GetUserGroupsResponse, error)
}

type GetUserGroupsUseCaseImpl struct {
	groupRepo repository.GroupRepository
	userID    entity.UserID
}

func NewGetUserGroupsUseCase(groupRepo repository.GroupRepository, userID entity.UserID) *GetUserGroupsUseCaseImpl {
	return &GetUserGroupsUseCaseImpl{
		groupRepo: groupRepo,
		userID:    userID,
	}
}

func (uc *GetUserGroupsUseCaseImpl) Execute() (*GetUserGroupsResponse, error) {
	groups, err := uc.groupRepo.FindGroupsByUserID(uc.userID)
	if err != nil {
		return nil, err
	}

	result := &GetUserGroupsResponse{
		Groups: make([]GroupInfo, 0, len(groups)),
	}

	for _, group := range groups {
		memberCount := len(group.GroupMembers)
		if group.GroupManagerID != "" {
			memberCount++
		}

		isCreator := group.GroupManagerID == uc.userID

		groupResponse := GroupInfo{
			ID:          string(group.GroupID),
			Name:        string(group.GroupName),
			MemberCount: memberCount,
			IsCreator:   isCreator,
		}

		result.Groups = append(result.Groups, groupResponse)
	}

	return result, nil
}
