package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
	"sync"
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
		return fmt.Errorf("not found group")
	}

	if groupFound.GroupManagerID == uc.userID {
		return fmt.Errorf("can not leave group, you are the manager")
	}

	memberFound := false
	var mutex sync.Mutex
	var wg sync.WaitGroup

	type Result struct {
		Index int
		Found bool
	}
	resultCh := make(chan Result, 1)

	numGoroutines := 4
	if len(groupFound.GroupMembers) < numGoroutines {
		numGoroutines = len(groupFound.GroupMembers)
	}

	chunkSize := (len(groupFound.GroupMembers) + numGoroutines - 1) / numGoroutines

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(startIdx int) {
			defer wg.Done()

			endIdx := startIdx + chunkSize
			if endIdx > len(groupFound.GroupMembers) {
				endIdx = len(groupFound.GroupMembers)
			}

			for i := startIdx; i < endIdx; i++ {
				if groupFound.GroupMembers[i] == uc.userID {
					select {
					case resultCh <- Result{Index: i, Found: true}:
					default:
					}
					return
				}
			}
		}(g * chunkSize)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	result, ok := <-resultCh
	if ok && result.Found {
		memberFound = true
		mutex.Lock()
		groupFound.GroupMembers = append(groupFound.GroupMembers[:result.Index], groupFound.GroupMembers[result.Index+1:]...)
		mutex.Unlock()
	}

	if !memberFound {
		return fmt.Errorf("user is not a member of the group")
	}

	_, err = uc.groupRepo.UpdateGroup(*groupFound)
	return err
}
