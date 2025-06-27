package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

type PostParticipationUseCase interface {
	Execute() (*entity.Event, error)
}
type PostParticipationUseCaseImpl struct {
	eventRepo repository.EventRepository
	groupRepo repository.GroupRepository
	userID    *entity.UserID
	eventID   *entity.EventID
	vote      *entity.Vote
}

func NewPostParticipationUseCase(eventRepo repository.EventRepository, groupRepo repository.GroupRepository, userID *entity.UserID, eventID *entity.EventID, vote *entity.Vote) *PostParticipationUseCaseImpl {
	return &PostParticipationUseCaseImpl{
		eventRepo: eventRepo,
		groupRepo: groupRepo,
		userID:    userID,
		eventID:   eventID,
		vote:      vote,
	}
}

func (uc *PostParticipationUseCaseImpl) Execute() (*entity.Event, error) {
	event, err := uc.eventRepo.FindEventByEventID(*uc.eventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	var isGroupMember bool = false

	groups, err := uc.groupRepo.FindGroupsByUserID(*uc.userID)
	if err != nil {
		return nil, fmt.Errorf("ユーザーの所属グループ取得中にエラーが発生しました: %v", err)
	}

	for _, group := range groups {
		for _, eventID := range group.GroupEvents {
			if eventID == event.EventID {
				isGroupMember = true
				break
			}
		}
		if isGroupMember {
			break
		}
	}

	if !isGroupMember {
		return nil, fmt.Errorf("not a group member. cannot vote")
	}

	found := false
	for i, member := range event.VotedMembers {
		if member.UserID == *uc.userID {
			event.VotedMembers[i] = entity.VotedMember{
				UserID: *uc.userID,
				Vote:   *uc.vote,
			}
			found = true
			break
		}
	}

	if !found {
		newMember := entity.VotedMember{
			UserID: *uc.userID,
			Vote:   *uc.vote,
		}
		event.VotedMembers = append(event.VotedMembers, newMember)
	}

	updatedEvent, err := uc.eventRepo.UpdateEvent(*event)
	if err != nil {
		return nil, fmt.Errorf("投票情報の更新に失敗しました: %v", err)
	}

	return updatedEvent, nil
}
