package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"sync"
)

type Member struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type GroupInfoResponse struct {
	GroupName string   `json:"group_name"`
	Password  string   `json:"password"`
	Members   []Member `json:"members"`
}

type GetGroupInfoUsecase struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
}

func NewGetGroupInfoUsecase(groupRepo repository.GroupRepository, userRepo repository.UserRepository) *GetGroupInfoUsecase {
	return &GetGroupInfoUsecase{
		groupRepo: groupRepo,
		userRepo:  userRepo,
	}
}

func (u *GetGroupInfoUsecase) Execute(groupID entity.GroupID) (*GroupInfoResponse, error) {
	groupName := entity.GroupName(groupID)
	group, err := u.groupRepo.FindGroupByGroupName(&groupName)
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

			user, err := u.userRepo.FindUserByUserID(id)
			if err != nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			members = append(members, Member{
				ID:   string(user.UserID),
				Name: string(user.UserName),
				Icon: string(user.UserIcon),
			})
		}(memberID)
	}

	wg.Wait()

	response := &GroupInfoResponse{
		GroupName: string(group.GroupName),
		Password:  string(group.GroupPassword),
		Members:   members,
	}

	return response, nil
}
