package repository_test

import (
	"chikokulympic-api/domain/entity"
	infra "chikokulympic-api/infrastructure/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestUserRepository(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	t.Run("FindUserByUserID", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			// テストケースのセットアップ
			expectedUser := &entity.User{
				UserID:   "user1",
				AuthID:   "auth1",
				UserName: "Test User",
				FCMToken: "fcm-token-1",
				Alias:    "tester",
			}

			// モックの応答を設定
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{
				{Key: "user_id", Value: expectedUser.UserID},
				{Key: "auth_id", Value: expectedUser.AuthID},
				{Key: "user_name", Value: expectedUser.UserName},
				{Key: "fcm_token", Value: expectedUser.FCMToken},
				{Key: "alias", Value: expectedUser.Alias},
			}))

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			user, err := userRepo.FindUserByUserID("user1")

			// 検証
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, expectedUser.UserID, user.UserID)
			assert.Equal(t, expectedUser.AuthID, user.AuthID)
			assert.Equal(t, expectedUser.UserName, user.UserName)
			assert.Equal(t, expectedUser.FCMToken, user.FCMToken)
			assert.Equal(t, expectedUser.Alias, user.Alias)
		})

		mt.Run("not found", func(mt *mtest.T) {
			// ドキュメントが見つからない場合のモック応答
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11000,
					Message: "not found",
					Name:    mongo.ErrNoDocuments.Error(),
				},
			))

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			user, err := userRepo.FindUserByUserID("nonexistent")

			// 検証
			assert.Nil(t, user)
			assert.Error(t, err)
		})
	})

	t.Run("FindUserByAuthID", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			// テストケースのセットアップ
			expectedUser := &entity.User{
				UserID:   "user1",
				AuthID:   "auth1",
				UserName: "Test User",
				FCMToken: "fcm-token-1",
				Alias:    "tester",
			}

			// モックの応答を設定
			mt.AddMockResponses(mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{
				{Key: "user_id", Value: expectedUser.UserID},
				{Key: "auth_id", Value: expectedUser.AuthID},
				{Key: "user_name", Value: expectedUser.UserName},
				{Key: "fcm_token", Value: expectedUser.FCMToken},
				{Key: "alias", Value: expectedUser.Alias},
			}))

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			user, err := userRepo.FindUserByAuthID("auth1")

			// 検証
			assert.NoError(t, err)
			assert.NotNil(t, user)
			assert.Equal(t, expectedUser.UserID, user.UserID)
			assert.Equal(t, expectedUser.AuthID, user.AuthID)
		})

		mt.Run("not found", func(mt *mtest.T) {
			// ドキュメントが見つからない場合のモック応答
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11000,
					Message: "not found",
					Name:    mongo.ErrNoDocuments.Error(),
				},
			))

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			user, err := userRepo.FindUserByAuthID("nonexistent")

			// 検証
			assert.Nil(t, user)
			assert.Error(t, err)
		})
	})

	t.Run("CreateUser", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			// テストケースのセットアップ
			newUser := &entity.User{
				UserID:   "user1",
				AuthID:   "auth1",
				UserName: "New User",
				FCMToken: "fcm-token-1",
				Alias:    "newbie",
			}

			// InsertOneのモック応答
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			createdUser, err := userRepo.CreateUser(newUser)

			// 検証
			assert.NoError(t, err)
			assert.NotNil(t, createdUser)
			assert.Equal(t, newUser.UserID, createdUser.UserID)
			assert.Equal(t, newUser.AuthID, createdUser.AuthID)
			assert.Equal(t, newUser.UserName, createdUser.UserName)
		})

		mt.Run("error", func(mt *mtest.T) {
			// 挿入時のエラーをシミュレート
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11000,
					Message: "duplicate key error",
				},
			))

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			newUser := &entity.User{
				UserID:   "user1",
				AuthID:   "auth1",
				UserName: "New User",
			}
			createdUser, err := userRepo.CreateUser(newUser)

			// 検証
			assert.Error(t, err)
			assert.Nil(t, createdUser)
		})
	})

	t.Run("UpdateUser", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			// テストケースのセットアップ
			userToUpdate := &entity.User{
				UserID:   "user1",
				AuthID:   "auth1",
				UserName: "Updated User",
				FCMToken: "fcm-token-updated",
				Alias:    "updated",
			}

			// UpdateOneのモック応答
			mt.AddMockResponses(
				mtest.CreateSuccessResponse(), // UpdateOne成功
				mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{ // FindOne成功
					{Key: "user_id", Value: userToUpdate.UserID},
					{Key: "auth_id", Value: userToUpdate.AuthID},
					{Key: "user_name", Value: userToUpdate.UserName},
					{Key: "fcm_token", Value: userToUpdate.FCMToken},
					{Key: "alias", Value: userToUpdate.Alias},
				}),
			)

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			updatedUser, err := userRepo.UpdateUser(userToUpdate)

			// 検証
			assert.NoError(t, err)
			assert.NotNil(t, updatedUser)
			assert.Equal(t, userToUpdate.UserID, updatedUser.UserID)
			assert.Equal(t, userToUpdate.UserName, updatedUser.UserName)
			assert.Equal(t, userToUpdate.FCMToken, updatedUser.FCMToken)
		})

		mt.Run("update error", func(mt *mtest.T) {
			// 更新時のエラーをシミュレート
			mt.AddMockResponses(mtest.CreateCommandErrorResponse(
				mtest.CommandError{
					Code:    11000,
					Message: "update failed",
				},
			))

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			userToUpdate := &entity.User{
				UserID:   "user1",
				UserName: "Updated User",
			}
			updatedUser, err := userRepo.UpdateUser(userToUpdate)

			// 検証
			assert.Error(t, err)
			assert.Nil(t, updatedUser)
		})
	})

	t.Run("DeleteUser", func(t *testing.T) {
		mt.Run("success", func(mt *mtest.T) {
			// テストケースのセットアップ
			userID := entity.UserID("user1")
			expectedUser := &entity.User{
				UserID:   userID,
				AuthID:   "auth1",
				UserName: "Test User",
				FCMToken: "fcm-token-1",
				Alias:    "tester",
			}

			// FindOneとDeleteOneのモック応答
			mt.AddMockResponses(
				mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{ // FindOne成功
					{Key: "user_id", Value: expectedUser.UserID},
					{Key: "auth_id", Value: expectedUser.AuthID},
					{Key: "user_name", Value: expectedUser.UserName},
					{Key: "fcm_token", Value: expectedUser.FCMToken},
					{Key: "alias", Value: expectedUser.Alias},
				}),
				mtest.CreateSuccessResponse(), // DeleteOne成功
			)

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			deletedUser, err := userRepo.DeleteUser(&userID)

			// 検証
			assert.NoError(t, err)
			assert.NotNil(t, deletedUser)
			assert.Equal(t, expectedUser.UserID, deletedUser.UserID)
		})

		mt.Run("not found", func(mt *mtest.T) {
			// ユーザーが見つからない場合のモック応答
			mt.AddMockResponses(
				mtest.CreateCommandErrorResponse(
					mtest.CommandError{
						Code:    11000,
						Message: "not found",
						Name:    mongo.ErrNoDocuments.Error(),
					},
				),
			)

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			userID := entity.UserID("nonexistent")
			deletedUser, err := userRepo.DeleteUser(&userID)

			// 検証
			assert.Error(t, err)
			assert.Nil(t, deletedUser)
		})

		mt.Run("delete error", func(mt *mtest.T) {
			// FindOneは成功するがDeleteOneが失敗するケース
			userID := entity.UserID("user1")
			expectedUser := &entity.User{
				UserID:   userID,
				AuthID:   "auth1",
				UserName: "Test User",
			}

			mt.AddMockResponses(
				mtest.CreateCursorResponse(1, "users", mtest.FirstBatch, bson.D{ // FindOne成功
					{Key: "user_id", Value: expectedUser.UserID},
					{Key: "auth_id", Value: expectedUser.AuthID},
					{Key: "user_name", Value: expectedUser.UserName},
				}),
				mtest.CreateCommandErrorResponse( // DeleteOne失敗
					mtest.CommandError{
						Code:    11000,
						Message: "delete failed",
					},
				),
			)

			// テスト対象のリポジトリを作成
			userRepo := infra.NewUserRepository(mt.DB)

			// テスト実行
			deletedUser, err := userRepo.DeleteUser(&userID)

			// 検証
			assert.Error(t, err)
			assert.Nil(t, deletedUser)
		})
	})
}
