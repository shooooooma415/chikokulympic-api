package repository_test

import (
	"context"
	"testing"

	"chikokulympic-api/domain/entity"
	"chikokulympic-api/infrastructure/mongo/repository"
	"chikokulympic-api/infrastructure/mongo/repository/testUtils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFindUserByUserID(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	testUser := &entity.User{
		UserID:   "test-user-id",
		AuthID:   "test-auth-id",
		UserName: "Test User",
		FCMToken: "fcm-token-123",
		Alias:    "tester",
	}

	_, err := db.Collection("users").InsertOne(context.Background(), testUser)
	assert.NoError(t, err)
	repo := repository.NewUserRepository(db)

	foundUser, err := repo.FindUserByUserID(testUser.UserID)

	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.UserID, foundUser.UserID)
	assert.Equal(t, testUser.AuthID, foundUser.AuthID)
	assert.Equal(t, testUser.UserName, foundUser.UserName)
}

func TestFindUserByAuthID(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	testUser := &entity.User{
		UserID:   "test-user-id",
		AuthID:   "test-auth-id",
		UserName: "Test User",
		FCMToken: "fcm-token-123",
		Alias:    "tester",
	}

	_, err := db.Collection("users").InsertOne(context.Background(), testUser)
	assert.NoError(t, err)

	repo := repository.NewUserRepository(db)

	foundUser, err := repo.FindUserByAuthID(testUser.AuthID)

	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.UserID, foundUser.UserID)
	assert.Equal(t, testUser.AuthID, foundUser.AuthID)
}

func TestCreateUser(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	testUser := &entity.User{
		UserID:   "new-user-id",
		AuthID:   "new-auth-id",
		UserName: "New User",
		FCMToken: "fcm-token-new",
		Alias:    "newbie",
	}

	createdUser, err := repo.CreateUser(testUser)

	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, testUser.UserID, createdUser.UserID)

	var savedUser entity.User
	err = db.Collection("users").FindOne(context.Background(), bson.M{"user_id": testUser.UserID}).Decode(&savedUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.UserName, savedUser.UserName)
}

func TestUpdateUser(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	testUser := &entity.User{
		UserID:   "update-user-id",
		AuthID:   "update-auth-id",
		UserName: "Update User",
		FCMToken: "fcm-token-update",
		Alias:    "updater",
	}

	_, err := db.Collection("users").InsertOne(context.Background(), testUser)
	assert.NoError(t, err)

	repo := repository.NewUserRepository(db)

	updatedUserName := entity.UserName("Updated User Name")
	updatedAlias := entity.Alias("updated-alias")
	testUser.UserName = updatedUserName
	testUser.Alias = updatedAlias

	updatedUser, err := repo.UpdateUser(testUser)

	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, updatedUserName, updatedUser.UserName)
	assert.Equal(t, updatedAlias, updatedUser.Alias)
}

func TestDeleteUser(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	testUser := &entity.User{
		UserID:   "delete-user-id",
		AuthID:   "delete-auth-id",
		UserName: "Delete User",
		FCMToken: "fcm-token-delete",
		Alias:    "deleter",
	}

	_, err := db.Collection("users").InsertOne(context.Background(), testUser)
	assert.NoError(t, err)

	repo := repository.NewUserRepository(db)

	userID := testUser.UserID
	deletedUser, err := repo.DeleteUser(&userID)

	assert.NoError(t, err)
	assert.NotNil(t, deletedUser)
	assert.Equal(t, testUser.UserID, deletedUser.UserID)

	var count int64
	count, err = db.Collection("users").CountDocuments(context.Background(), bson.M{"user_id": userID})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestFindUserNotFound(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	repo := repository.NewUserRepository(db)

	nonExistentID := entity.UserID("non-existent-id")
	user, err := repo.FindUserByUserID(nonExistentID)

	assert.NoError(t, err)
	assert.Nil(t, user)
}
