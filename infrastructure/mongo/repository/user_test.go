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

func TestUserRepository(t *testing.T) {
	// 各テストで共通のセットアップ処理
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()
	repo := repository.NewUserRepository(db)

	t.Run("FindUserByUserID", func(t *testing.T) {
		testCases := []struct {
			name     string
			user     *entity.User
			userID   entity.UserID
			expected *entity.User
			isFound  bool
		}{
			{
				name: "正常系: 存在するユーザーIDで検索",
				user: &entity.User{
					UserID:   "test-user-id-1",
					AuthID:   "test-auth-id-1",
					UserName: "Test User 1",
					FCMToken: "fcm-token-123",
					Alias:    "tester1",
				},
				userID:   "test-user-id-1",
				expected: nil, // 後でセット
				isFound:  true,
			},
			{
				name:     "異常系: 存在しないユーザーIDで検索",
				user:     nil,
				userID:   "non-existent-id",
				expected: nil,
				isFound:  false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				if tc.user != nil {
					_, err := db.Collection("users").InsertOne(context.Background(), tc.user)
					assert.NoError(t, err)
					tc.expected = tc.user
				}

				// テスト実行
				foundUser, err := repo.FindUserByUserID(tc.userID)

				// 結果の検証
				assert.NoError(t, err)
				if tc.isFound {
					assert.NotNil(t, foundUser)
					assert.Equal(t, tc.expected.UserID, foundUser.UserID)
					assert.Equal(t, tc.expected.AuthID, foundUser.AuthID)
					assert.Equal(t, tc.expected.UserName, foundUser.UserName)
				} else {
					assert.Nil(t, foundUser)
				}

				// クリーンアップ
				if tc.user != nil {
					_, err = db.Collection("users").DeleteMany(context.Background(), bson.M{"user_id": tc.user.UserID})
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("FindUserByAuthID", func(t *testing.T) {
		testCases := []struct {
			name     string
			user     *entity.User
			authID   string
			expected *entity.User
			isFound  bool
		}{
			{
				name: "正常系: 存在する認証IDで検索",
				user: &entity.User{
					UserID:   "test-user-id-2",
					AuthID:   "test-auth-id-2",
					UserName: "Test User 2",
					FCMToken: "fcm-token-456",
					Alias:    "tester2",
				},
				authID:   "test-auth-id-2",
				expected: nil, // 後でセット
				isFound:  true,
			},
			{
				name:     "異常系: 存在しない認証IDで検索",
				user:     nil,
				authID:   "non-existent-auth-id",
				expected: nil,
				isFound:  false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				if tc.user != nil {
					_, err := db.Collection("users").InsertOne(context.Background(), tc.user)
					assert.NoError(t, err)
					tc.expected = tc.user
				}

				// テスト実行
				foundUser, err := repo.FindUserByAuthID(entity.AuthID(tc.authID))

				// 結果の検証
				assert.NoError(t, err)
				if tc.isFound {
					assert.NotNil(t, foundUser)
					assert.Equal(t, tc.expected.UserID, foundUser.UserID)
					assert.Equal(t, tc.expected.AuthID, foundUser.AuthID)
				} else {
					assert.Nil(t, foundUser)
				}

				// クリーンアップ
				if tc.user != nil {
					_, err = db.Collection("users").DeleteMany(context.Background(), bson.M{"auth_id": tc.user.AuthID})
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("CreateUser", func(t *testing.T) {
		testCases := []struct {
			name  string
			user  *entity.User
			error bool
		}{
			{
				name: "正常系: 新規ユーザー作成",
				user: &entity.User{
					// UserIDはリポジトリで自動生成されるので設定しない
					AuthID:   "new-auth-id",
					UserName: "New User",
					FCMToken: "fcm-token-new",
					Alias:    "newbie",
				},
				error: false,
			},
			{
				name: "正常系: IDを事前に指定して新規ユーザー作成",
				user: &entity.User{
					UserID:   "507f1f77bcf86cd799439011", // 有効なObjectIDの形式
					AuthID:   "new-auth-id-2",
					UserName: "New User 2",
					FCMToken: "fcm-token-new-2",
					Alias:    "newbie2",
				},
				error: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テスト実行
				createdUser, err := repo.CreateUser(*tc.user)

				// 結果の検証
				if tc.error {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, createdUser)

					if tc.user.UserID != "" {
						// IDが事前に設定されている場合は同じIDが使われていることを確認
						assert.Equal(t, tc.user.UserID, createdUser.UserID)
					} else {
						// 自動生成された場合はIDが空でないことを確認
						assert.NotEmpty(t, createdUser.UserID)
					}

					// DBに保存されていることを確認
					var savedUser entity.User
					var err error
					if tc.user.UserID != "" {
						err = db.Collection("users").FindOne(context.Background(), bson.M{"user_id": tc.user.UserID}).Decode(&savedUser)
					} else {
						err = db.Collection("users").FindOne(context.Background(), bson.M{"user_id": createdUser.UserID}).Decode(&savedUser)
					}
					assert.NoError(t, err)
					assert.Equal(t, tc.user.UserName, savedUser.UserName)
				}

				// クリーンアップ
				if tc.user.UserID != "" {
					_, err = db.Collection("users").DeleteMany(context.Background(), bson.M{"user_id": tc.user.UserID})
				} else {
					_, err = db.Collection("users").DeleteMany(context.Background(), bson.M{"user_id": createdUser.UserID})
				}
				assert.NoError(t, err)
			})
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		testCases := []struct {
			name        string
			initialUser *entity.User
			updatedUser *entity.User
			error       bool
		}{
			{
				name: "正常系: ユーザー情報更新",
				initialUser: &entity.User{
					UserID:   "update-user-id",
					AuthID:   "update-auth-id",
					UserName: "Update User",
					FCMToken: "fcm-token-update",
					Alias:    "updater",
				},
				updatedUser: &entity.User{
					UserID:   "update-user-id",
					AuthID:   "update-auth-id",
					UserName: "Updated User Name",
					FCMToken: "fcm-token-update",
					Alias:    "updated-alias",
				},
				error: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				_, err := db.Collection("users").InsertOne(context.Background(), tc.initialUser)
				assert.NoError(t, err)

				// テスト実行
				updatedUser, err := repo.UpdateUser(*tc.updatedUser)

				// 結果の検証
				if tc.error {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, updatedUser)
					assert.Equal(t, tc.updatedUser.UserName, updatedUser.UserName)
					assert.Equal(t, tc.updatedUser.Alias, updatedUser.Alias)

					// DBが更新されたことを確認
					var savedUser entity.User
					err = db.Collection("users").FindOne(context.Background(), bson.M{"user_id": tc.initialUser.UserID}).Decode(&savedUser)
					assert.NoError(t, err)
					assert.Equal(t, tc.updatedUser.UserName, savedUser.UserName)
					assert.Equal(t, tc.updatedUser.Alias, savedUser.Alias)
				}

				// クリーンアップ
				_, err = db.Collection("users").DeleteMany(context.Background(), bson.M{"user_id": tc.initialUser.UserID})
				assert.NoError(t, err)
			})
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		testCases := []struct {
			name  string
			user  *entity.User
			error bool
		}{
			{
				name: "正常系: ユーザー削除",
				user: &entity.User{
					UserID:   "delete-user-id",
					AuthID:   "delete-auth-id",
					UserName: "Delete User",
					FCMToken: "fcm-token-delete",
					Alias:    "deleter",
				},
				error: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				_, err := db.Collection("users").InsertOne(context.Background(), tc.user)
				assert.NoError(t, err)

				// テスト実行
				deletedUser, err := repo.DeleteUser(*tc.user)

				// 結果の検証
				if tc.error {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, deletedUser)
					assert.Equal(t, tc.user.UserID, deletedUser.UserID)

					// DBから削除されたことを確認
					var count int64
					count, err = db.Collection("users").CountDocuments(context.Background(), bson.M{"user_id": tc.user.UserID})
					assert.NoError(t, err)
					assert.Equal(t, int64(0), count)
				}
			})
		}
	})
}
