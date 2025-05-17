package repository

import (
	"context"
	"fmt"
	"time"

	"chikokulympic-api/domain/entity"
	repo "chikokulympic-api/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GroupRepo struct {
	groupCollection *mongo.Collection
}

func NewGroupRepository(db *mongo.Database) repo.GroupRepository {
	return &GroupRepo{
		groupCollection: db.Collection("groups"),
	}
}

func (gr *GroupRepo) FindGroupByGroupName(groupName *entity.GroupName) (*entity.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var group entity.Group
	filter := bson.M{"name": string(*groupName)} // ポインタをstring型に変換
	err := gr.groupCollection.FindOne(ctx, filter).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("group not found with name: %s", string(*groupName))
		}
		return nil, fmt.Errorf("error finding group by name: %w", err)
	}

	return &group, nil
}

func (gr *GroupRepo) FindGroupsByUserID(userID entity.UserID) ([]*entity.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var groups []*entity.Group
	filter := bson.M{
		"$or": []bson.M{
			{"manager_id": userID},
			{"members": bson.M{"$in": []entity.UserID{userID}}},
		},
	}

	cursor, err := gr.groupCollection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error finding groups by user ID: %w", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &groups); err != nil {
		return nil, fmt.Errorf("error decoding groups: %w", err)
	}

	if len(groups) == 0 {
		return []*entity.Group{}, nil // 空の配列を返す（nilではなく）
	}

	return groups, nil
}

func (gr *GroupRepo) CreateGroup(group *entity.Group) (*entity.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := gr.groupCollection.InsertOne(ctx, group)
	if err != nil {
		return nil, fmt.Errorf("error creating group: %w", err)
	}

	return group, nil
}

func (gr *GroupRepo) UpdateGroup(group *entity.Group) (*entity.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"group_id": group.GroupID}
	update := bson.M{"$set": group}

	_, err := gr.groupCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("error updating group: %w", err)
	}

	return group, nil
}

func (gr *GroupRepo) DeleteGroup(group *entity.Group) (*entity.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"group_id": group.GroupID}

	_, err := gr.groupCollection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error deleting group: %w", err)
	}

	return group, nil
}

func (gr *GroupRepo) FindGroupByGroupID(groupID entity.GroupID) (*entity.Group, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var group entity.Group
	filter := bson.M{"group_id": groupID}
	err := gr.groupCollection.FindOne(ctx, filter).Decode(&group)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("group not found with ID: %s", string(groupID))
		}
		return nil, fmt.Errorf("error finding group by ID: %w", err)
	}

	return &group, nil
}
