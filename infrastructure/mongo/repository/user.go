package repository

import (
	"context"
	"errors"
	"fmt"

	"chikokulympic-api/domain/entity"
	repo "chikokulympic-api/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) repo.UserRepository {
	return &userRepository{
		userCollection: db.Collection("users"),
	}
}

func (r *userRepository) FindUserByUserID(userID entity.UserID) (*entity.User, error) {
	var user entity.User
	err := r.userCollection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindUserByAuthID(authID entity.AuthID) (*entity.User, error) {
	var user entity.User
	err := r.userCollection.FindOne(context.Background(), bson.M{"auth_id": authID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user entity.User) (*entity.User, error) {
	result, err := r.userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}

	if result.InsertedID == nil {
		return nil, fmt.Errorf("failed to insert user")
	}

	return &user, nil
}

func (r *userRepository) DeleteUser(user entity.User) (*entity.User, error) {
	var deletedUser entity.User
	filter := bson.M{"user_id": user.UserID}

	err := r.userCollection.FindOneAndDelete(context.Background(), filter).Decode(&deletedUser)
	if err != nil {
		return nil, err
	}

	return &deletedUser, nil
}

func (r *userRepository) UpdateUser(user entity.User) (*entity.User, error) {
	filter := bson.M{"user_id": user.UserID}
	update := bson.M{"$set": user}

	_, err := r.userCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	updatedUser, err := r.FindUserByUserID(user.UserID)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
