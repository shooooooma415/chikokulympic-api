package repository

import (
	"context"
	"testing"

	"chikokulympic-api/domain/entity"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// テスト用のMongoDBクライアントとデータベース接続を設定する
func setupTestDB(t *testing.T) (*mongo.Database, func()) {
	// テスト用のMongoDBに接続（ローカルかテスト用のコンテナを想定）
	// 実際の環境に合わせて適宜修正してください
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// テスト用の一時データベースを作成
	db := client.Database("test_chikokulympic")

	// クリーンアップ関数を返す
	cleanup := func() {
		err := db.Drop(context.Background())
		if err != nil {
			t.Fatalf("Failed to drop test database: %v", err)
		}
		client.Disconnect(context.Background())
	}

	return db, cleanup
}

func TestFindUserByUserID(t *testing.T) {
	// テスト用DBをセットアップ
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// テスト用ユーザーを作成
	testUser := &entity.User{
		UserID:   "test-user-id",
		AuthID:   "test-auth-id",
		UserName: "Test User",
		FCMToken: "fcm-token-123",
		Alias:    "tester",
	}

	// テストデータを直接コレクションに挿入
	_, err := db.Collection("users").InsertOne(context.Background(), testUser)
	assert.NoError(t, err)

	// テスト対象のリポジトリを作成
	repo := NewUserRepository(db)

	// テスト実行
	foundUser, err := repo.FindUserByUserID(testUser.UserID)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.UserID, foundUser.UserID)
	assert.Equal(t, testUser.AuthID, foundUser.AuthID)
	assert.Equal(t, testUser.UserName, foundUser.UserName)
}

func TestFindUserByAuthID(t *testing.T) {
	// テスト用DBをセットアップ
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// テスト用ユーザーを作成
	testUser := &entity.User{
		UserID:   "test-user-id",
		AuthID:   "test-auth-id",
		UserName: "Test User",
		FCMToken: "fcm-token-123",
		Alias:    "tester",
	}

	// テストデータを直接コレクションに挿入
	_, err := db.Collection("users").InsertOne(context.Background(), testUser)
	assert.NoError(t, err)

	// テスト対象のリポジトリを作成
	repo := NewUserRepository(db)

	// テスト実行
	foundUser, err := repo.FindUserByAuthID(testUser.AuthID)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, testUser.UserID, foundUser.UserID)
	assert.Equal(t, testUser.AuthID, foundUser.AuthID)
}

func TestCreateUser(t *testing.T) {
	// テスト用DBをセットアップ
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// テスト対象のリポジトリを作成
	repo := NewUserRepository(db)

	// テスト用ユーザーを作成
	testUser := &entity.User{
		UserID:   "new-user-id",
		AuthID:   "new-auth-id",
		UserName: "New User",
		FCMToken: "fcm-token-new",
		Alias:    "newbie",
	}

	// テスト実行
	createdUser, err := repo.CreateUser(testUser)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, createdUser)
	assert.Equal(t, testUser.UserID, createdUser.UserID)

	// データベースに保存されたことを確認
	var savedUser entity.User
	err = db.Collection("users").FindOne(context.Background(), bson.M{"user_id": testUser.UserID}).Decode(&savedUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.UserName, savedUser.UserName)
}

func TestUpdateUser(t *testing.T) {
	// テスト用DBをセットアップ
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// テスト用ユーザーを作成
	testUser := &entity.User{
		UserID:   "update-user-id",
		AuthID:   "update-auth-id",
		UserName: "Update User",
		FCMToken: "fcm-token-update",
		Alias:    "updater",
	}

	// テストデータを直接コレクションに挿入
	_, err := db.Collection("users").InsertOne(context.Background(), testUser)
	assert.NoError(t, err)

	// テスト対象のリポジトリを作成
	repo := NewUserRepository(db)

	// 更新用データを準備
	updatedUserName := entity.UserName("Updated User Name")
	updatedAlias := entity.Alias("updated-alias")
	testUser.UserName = updatedUserName
	testUser.Alias = updatedAlias

	// テスト実行
	updatedUser, err := repo.UpdateUser(testUser)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser)
	assert.Equal(t, updatedUserName, updatedUser.UserName)
	assert.Equal(t, updatedAlias, updatedUser.Alias)
}

func TestDeleteUser(t *testing.T) {
	// テスト用DBをセットアップ
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// テスト用ユーザーを作成
	testUser := &entity.User{
		UserID:   "delete-user-id",
		AuthID:   "delete-auth-id",
		UserName: "Delete User",
		FCMToken: "fcm-token-delete",
		Alias:    "deleter",
	}

	// テストデータを直接コレクションに挿入
	_, err := db.Collection("users").InsertOne(context.Background(), testUser)
	assert.NoError(t, err)

	// テスト対象のリポジトリを作成
	repo := NewUserRepository(db)

	// テスト実行
	userID := testUser.UserID
	deletedUser, err := repo.DeleteUser(&userID)

	// 検証
	assert.NoError(t, err)
	assert.NotNil(t, deletedUser)
	assert.Equal(t, testUser.UserID, deletedUser.UserID)

	// 実際に削除されたか確認
	var count int64
	count, err = db.Collection("users").CountDocuments(context.Background(), bson.M{"user_id": userID})
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)
}

func TestFindUserNotFound(t *testing.T) {
	// テスト用DBをセットアップ
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// テスト対象のリポジトリを作成
	repo := NewUserRepository(db)

	// 存在しないユーザーIDで検索
	nonExistentID := entity.UserID("non-existent-id")
	user, err := repo.FindUserByUserID(nonExistentID)

	// 検証（エラーはなくユーザーはnilが返る仕様）
	assert.NoError(t, err)
	assert.Nil(t, user)
}
