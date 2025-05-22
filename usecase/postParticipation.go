package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
	"time"
)

type PostParticipationUseCase interface {
	Execute() (*entity.Event, error)
}
type PostParticipationUseCaseImpl struct {
	eventRepo repository.EventRepository
	userID    *entity.UserID
	eventID   *entity.EventID
	vote      *entity.Vote
}

func NewPostParticipationUseCase(eventRepo repository.EventRepository, userID *entity.UserID, eventID *entity.EventID, vote *entity.Vote) *PostParticipationUseCaseImpl {
	return &PostParticipationUseCaseImpl{
		eventRepo: eventRepo,
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

	// voteの値を参加/不参加として処理
	voteValue := string(*uc.vote)

	// 既に同じユーザーのエントリがあるか確認
	found := false
	for i, member := range event.VotedMembers {
		if member.UserID == *uc.userID {
			// 既存のメンバーの投票情報を更新
			// 参加ステータスは更新しない(isArrivalは変更しない)
			// ArrivalDateTimeはvoteが"yes"または"no"の場合に更新
			if voteValue == "yes" || voteValue == "no" {
				// 現在時刻を設定
				event.VotedMembers[i].ArrivalDateTime = time.Now()
			}
			found = true
			break
		}
	}

	// 新しいメンバーの場合は追加
	if !found {
		// 初期状態ではUserIDのみ設定
		newMember := entity.VotedMember{
			UserID: *uc.userID,
		}

		// voteが"yes"または"no"の場合
		if voteValue == "yes" || voteValue == "no" {
			// ArrivalDateTimeを設定
			newMember.ArrivalDateTime = time.Now()
		}

		event.VotedMembers = append(event.VotedMembers, newMember)
	}

	// イベント情報を更新
	updatedEvent, err := uc.eventRepo.UpdateEvent(*event)
	if err != nil {
		return nil, fmt.Errorf("参加情報の更新に失敗しました: %v", err)
	}

	return updatedEvent, nil
}
