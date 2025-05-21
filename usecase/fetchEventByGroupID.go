package usecase

import (
	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"
	"fmt"
)

type FetchEventInfoUsecase interface {
	Execute(groupID entity.EventID) (*entity.Group, error)
}
type FetchGroupInfoUsecaseImpl struct {
	groupRepo repository.GroupRepository
	groupID   *entity.GroupID
}

func NewFetchEventInfoUsecase(groupRepo repository.GroupRepository, groupID *entity.GroupID) *FetchGroupInfoUsecaseImpl {
	return &FetchGroupInfoUsecaseImpl{
		groupRepo: groupRepo,
		groupID:   groupID,
	}
}

func (uc *FetchGroupInfoUsecaseImpl) Execute() (*entity.Group, error) {
	group, err := uc.groupRepo.FindGroupByGroupID(*uc.groupID)
	if err != nil {
		return nil, fmt.Errorf("グループ情報の取得中にエラーが発生しました: %v", err)
	}
	if group == nil {
		return nil, fmt.Errorf("指定されたグループID %s が存在しません", *uc.groupID)
	}
	return group, nil
}
