package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"sync"
)


type GroupResponse struct {
	ID          string `json:"id" example:"group123"`
	Name        string `json:"name" example:"テストグループ"`
	MemberCount int    `json:"member_count" example:"5"`
	IsCreator   bool   `json:"is_creator" example:"true"`
}


type FetchUserGroupsResponse struct {
	Groups []GroupResponse `json:"groups"`
}


type UserGroup FetchUserGroupsResponse

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

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, group := range groups {
		wg.Add(1)
		go func(g *entity.Group) {
			defer wg.Done()

			memberCount := len(g.GroupMembers)
			isCreator := g.GroupManagerID == uc.userID

			groupResponse := GroupResponse{
				ID:          string(g.GroupID),
				Name:        string(g.GroupName),
				MemberCount: memberCount,
				IsCreator:   isCreator,
			}

			mu.Lock()
			result.Groups = append(result.Groups, groupResponse)
			mu.Unlock()
		}(group)
	}

	wg.Wait()

	return result, nil
}
