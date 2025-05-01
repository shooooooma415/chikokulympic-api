package repository

import (
	"context"
	"errors"

	"chikokulympic-api/domain/entity"
	"chikokulympic-api/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

// NewUserRepository は新しいUserRepositoryのインスタンスを作成します
func NewUserRepository(db *mongo.Database) repository.UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

// FindUserByUserID はユーザーIDでユーザーを検索します
func (r *userRepository) FindUserByUserID(userID entity.UserID) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // ユーザーが見つからない場合
		}
		return nil, err
	}
	return &user, nil
}

// FindUserByAuthID は認証IDでユーザーを検索します
func (r *userRepository) FindUserByAuthID(authID entity.AuthID) (*entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(context.Background(), bson.M{"auth_id": authID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // ユーザーが見つからない場合
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser は新しいユーザーを作成します
func (r *userRepository) CreateUser(user *entity.User) (*entity.User, error) {
	_, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// DeleteUser はユーザーIDでユーザーを削除します
func (r *userRepository) DeleteUser(userID *entity.UserID) (*entity.User, error) {
	// 削除する前にユーザーを取得
	user, err := r.FindUserByUserID(*userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// ユーザーを削除
	_, err = r.collection.DeleteOne(context.Background(), bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser はユーザー情報を更新します
func (r *userRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	filter := bson.M{"user_id": user.UserID}
	update := bson.M{"$set": user}

	_, err := r.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	// 更新後のユーザー情報を取得
	updatedUser, err := r.FindUserByUserID(user.UserID)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
