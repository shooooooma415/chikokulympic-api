package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

type LeaveGroupUseCase interface {
	Execute() error
}

type LeaveGroupUseCaseImpl struct {
	groupRepo repository.GroupRepository
	userID    entity.UserID
	groupID   entity.GroupID
}

func NewLeaveGroupUseCase(groupRepo repository.GroupRepository, userID entity.UserID, groupID entity.GroupID) *LeaveGroupUseCaseImpl {
	return &LeaveGroupUseCaseImpl{
		groupRepo: groupRepo,
		userID:    userID,
		groupID:   groupID,
	}
}

func (uc *LeaveGroupUseCaseImpl) Execute() error {
	groupFound, err := uc.groupRepo.FindGroupByGroupID(uc.groupID)
	if err != nil {
		return err
	}
	if groupFound == nil {
		return fmt.Errorf("グループが見つかりません")
	}

	for i, memberID := range groupFound.GroupMembers {
		if memberID == uc.userID {
			groupFound.GroupMembers = append(groupFound.GroupMembers[:i], groupFound.GroupMembers[i+1:]...)
			break
		}
	}

	_, err = uc.groupRepo.UpdateGroup(groupFound)
	return err
}
