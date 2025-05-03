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

func TestFindGroupByGroupName(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	testGroup := &entity.Group{
		GroupID:          "test-group-id",
		GroupName:        "TestGroup",
		GroupPassword:    "password123",
		GroupManagerID:   "manager-user-id",
		GroupDescription: "Test group description",
		GroupMembers:     []entity.UserID{"member1-id", "member2-id"},
		GroupEvents:      []entity.EventID{"event1-id", "event2-id"},
	}

	_, err := db.Collection("groups").InsertOne(context.Background(), testGroup)
	assert.NoError(t, err)

	repo := repository.NewGroupRepo(db)

	foundGroup, err := repo.FindGroupByGroupName(&testGroup.GroupName)

	assert.NoError(t, err)
	assert.NotNil(t, foundGroup)
	assert.Equal(t, testGroup.GroupID, foundGroup.GroupID)
	assert.Equal(t, testGroup.GroupName, foundGroup.GroupName)
	assert.Equal(t, testGroup.GroupManagerID, foundGroup.GroupManagerID)
}

func TestFindGroupByUserID(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	managerID := entity.UserID("manager-user-id")
	memberID := entity.UserID("member1-id")

	testGroup := &entity.Group{
		GroupID:          "test-group-id",
		GroupName:        "TestGroup",
		GroupPassword:    "password123",
		GroupManagerID:   managerID,
		GroupDescription: "Test group description",
		GroupMembers:     []entity.UserID{memberID, "member2-id"},
		GroupEvents:      []entity.EventID{"event1-id", "event2-id"},
	}

	_, err := db.Collection("groups").InsertOne(context.Background(), testGroup)
	assert.NoError(t, err)

	repo := repository.NewGroupRepo(db)

	// マネージャーIDでの検索テスト
	foundGroup, err := repo.FindGroupByUserID(managerID)
	assert.NoError(t, err)
	assert.NotNil(t, foundGroup)
	assert.Equal(t, testGroup.GroupID, foundGroup.GroupID)

	// メンバーIDでの検索テスト
	foundGroup, err = repo.FindGroupByUserID(memberID)
	assert.NoError(t, err)
	assert.NotNil(t, foundGroup)
	assert.Equal(t, testGroup.GroupID, foundGroup.GroupID)
}

func TestCreateGroup(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	repo := repository.NewGroupRepo(db)

	testGroup := &entity.Group{
		GroupID:          "new-group-id",
		GroupName:        "NewGroup",
		GroupPassword:    "newpassword",
		GroupManagerID:   "new-manager-id",
		GroupDescription: "New group description",
		GroupMembers:     []entity.UserID{"new-member1-id", "new-member2-id"},
		GroupEvents:      []entity.EventID{"new-event1-id", "new-event2-id"},
	}

	createdGroup, err := repo.CreateGroup(testGroup)

	assert.NoError(t, err)
	assert.NotNil(t, createdGroup)
	assert.Equal(t, testGroup.GroupID, createdGroup.GroupID)
	assert.Equal(t, testGroup.GroupName, createdGroup.GroupName)

	var savedGroup entity.Group
	err = db.Collection("groups").FindOne(context.Background(), bson.M{"group_id": testGroup.GroupID}).Decode(&savedGroup)
	assert.NoError(t, err)
	assert.Equal(t, testGroup.GroupName, savedGroup.GroupName)
	assert.Equal(t, testGroup.GroupDescription, savedGroup.GroupDescription)
}

func TestCreateGroupDuplicate(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	repo := repository.NewGroupRepo(db)

	testGroup := &entity.Group{
		GroupID:          "existing-group-id",
		GroupName:        "ExistingGroup",
		GroupPassword:    "password123",
		GroupManagerID:   "manager-id",
		GroupDescription: "Existing group description",
		GroupMembers:     []entity.UserID{"member1-id", "member2-id"},
		GroupEvents:      []entity.EventID{"event1-id", "event2-id"},
	}

	// 最初のグループ作成
	_, err := repo.CreateGroup(testGroup)
	assert.NoError(t, err)

	// 同じ名前で再度作成を試みる
	duplicateGroup := &entity.Group{
		GroupID:          "another-group-id",
		GroupName:        testGroup.GroupName, // 同じ名前
		GroupPassword:    "anotherpassword",
		GroupManagerID:   "another-manager-id",
		GroupDescription: "Another group description",
	}

	_, err = repo.CreateGroup(duplicateGroup)
	assert.Error(t, err) // エラーが返されることを期待
}

func TestUpdateGroup(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	testGroup := &entity.Group{
		GroupID:          "update-group-id",
		GroupName:        "UpdateGroup",
		GroupPassword:    "updatepassword",
		GroupManagerID:   "update-manager-id",
		GroupDescription: "Update group description",
		GroupMembers:     []entity.UserID{"update-member1-id", "update-member2-id"},
		GroupEvents:      []entity.EventID{"update-event1-id", "update-event2-id"},
	}

	_, err := db.Collection("groups").InsertOne(context.Background(), testGroup)
	assert.NoError(t, err)

	repo := repository.NewGroupRepo(db)

	// グループ情報の更新
	updatedDescription := entity.GroupDescription("Updated group description")
	updatedMembers := entity.GroupMembers{"update-member1-id", "update-member2-id", "update-member3-id"}

	testGroup.GroupDescription = updatedDescription
	testGroup.GroupMembers = updatedMembers

	updatedGroup, err := repo.UpdateGroup(testGroup)

	assert.NoError(t, err)
	assert.NotNil(t, updatedGroup)
	assert.Equal(t, updatedDescription, updatedGroup.GroupDescription)
	assert.Equal(t, len(updatedMembers), len(updatedGroup.GroupMembers))

	// DBから直接取得して確認
	var savedGroup entity.Group
	err = db.Collection("groups").FindOne(context.Background(), bson.M{"group_id": testGroup.GroupID}).Decode(&savedGroup)
	assert.NoError(t, err)
	assert.Equal(t, updatedDescription, savedGroup.GroupDescription)
	assert.Equal(t, len(updatedMembers), len(savedGroup.GroupMembers))
}

func TestDeleteGroup(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	testGroup := &entity.Group{
		GroupID:          "delete-group-id",
		GroupName:        "DeleteGroup",
		GroupPassword:    "deletepassword",
		GroupManagerID:   "delete-manager-id",
		GroupDescription: "Delete group description",
		GroupMembers:     []entity.UserID{"delete-member1-id", "delete-member2-id"},
		GroupEvents:      []entity.EventID{"delete-event1-id", "delete-event2-id"},
	}

	_, err := db.Collection("groups").InsertOne(context.Background(), testGroup)
	assert.NoError(t, err)

	repo := repository.NewGroupRepo(db)

	deletedGroup, err := repo.DeleteGroup(testGroup)

	assert.NoError(t, err)
	assert.NotNil(t, deletedGroup)
	assert.Equal(t, testGroup.GroupID, deletedGroup.GroupID)

	var count int64
	count, err = db.Collection("groups").CountDocuments(context.Background(), bson.M{"group_id": testGroup.GroupID})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestFindGroupNotFound(t *testing.T) {
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()

	repo := repository.NewGroupRepo(db)

	nonExistentName := entity.GroupName("NonExistentGroup")
	group, err := repo.FindGroupByGroupName(&nonExistentName)

	assert.Error(t, err)
	assert.Nil(t, group)

	nonExistentUserID := entity.UserID("non-existent-user-id")
	group, err = repo.FindGroupByUserID(nonExistentUserID)

	assert.Error(t, err)
	assert.Nil(t, group)
}
