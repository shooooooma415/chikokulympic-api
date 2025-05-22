package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"sync"
)

type FetchGroupInfoUsecase interface {
	Execute() (*GroupInfoResponse, error)
}

type GroupInfoFetcherUsecaseImpl struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
	groupID   *entity.GroupID
}

type Member struct {
	ID   entity.UserID   `json:"id"`
	Name entity.UserName `json:"name"`
	Icon entity.UserIcon `json:"icon"`
}

type GroupInfoResponse struct {
	GroupName      entity.GroupName     `json:"group_name"`
	Password       entity.GroupPassword `json:"password"`
	Members        []Member             `json:"members"`
	GroupManagerID entity.UserID        `json:"group_manager_id"`
}


func NewFetchGroupInfoUsecase(groupRepo repository.GroupRepository, userRepo repository.UserRepository, groupID *entity.GroupID) FetchGroupInfoUsecase {
	return &GroupInfoFetcherUsecaseImpl{
		groupRepo: groupRepo,
		userRepo:  userRepo,
		groupID:   groupID,
	}
}

func (uc *GroupInfoFetcherUsecaseImpl) Execute() (*GroupInfoResponse, error) {
	group, err := uc.groupRepo.FindGroupByGroupID(*uc.groupID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	members := make([]Member, 0, len(group.GroupMembers))

	for _, memberID := range group.GroupMembers {
		wg.Add(1)
		go func(id entity.UserID) {
			defer wg.Done()

			user, err := uc.userRepo.FindUserByUserID(id)
			if err != nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			members = append(members, Member{
				ID:   user.UserID,
				Name: user.UserName,
				Icon: user.UserIcon,
			})
		}(memberID)
	}

	wg.Wait()

	response := &GroupInfoResponse{
		GroupName:      group.GroupName,
		Password:       group.GroupPassword,
		Members:        members,
		GroupManagerID: group.GroupManagerID,
	}

	return response, nil
}
