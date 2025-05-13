package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
)

type GroupResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MemberCount int    `json:"member_count"`
	IsCreator   bool   `json:"is_creator"`
}

type FetchUserGroupsResponse struct {
	Groups []GroupResponse `json:"groups"`
}

type FetchUserGroupsUseCase interface {
	Execute() (*FetchUserGroupsResponse, error)
}

type FetchUserGroupsUseCaseImpl struct {
	groupRepo repository.GroupRepository
	userID    entity.UserID
}

func NewFetchUserGroupsUseCase(groupRepo repository.GroupRepository, userID entity.UserID) *FetchUserGroupsUseCaseImpl {
	return &FetchUserGroupsUseCaseImpl{
		groupRepo: groupRepo,
		userID:    userID,
	}
}

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
