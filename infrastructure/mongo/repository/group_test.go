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

func TestGroupRepository(t *testing.T) {
	// 各テストで共通のセットアップ処理
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()
	repo := repository.NewGroupRepo(db)

	t.Run("FindGroupByGroupName", func(t *testing.T) {
		testCases := []struct {
			name        string
			group       *entity.Group
			groupName   entity.GroupName
			expected    *entity.Group
			shouldError bool
		}{
			{
				name: "正常系: 存在するグループ名で検索",
				group: &entity.Group{
					GroupID:          "test-group-id-1",
					GroupName:        "TestGroup1",
					GroupPassword:    "password123",
					GroupManagerID:   "manager-user-id-1",
					GroupDescription: "Test group description 1",
					GroupMembers:     []entity.UserID{"member1-id", "member2-id"},
					GroupEvents:      []entity.EventID{"event1-id", "event2-id"},
				},
				groupName:   "TestGroup1",
				expected:    nil, // 後でセット
				shouldError: false,
			},
			{
				name:        "異常系: 存在しないグループ名で検索",
				group:       nil,
				groupName:   "NonExistentGroup",
				expected:    nil,
				shouldError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				if tc.group != nil {
					_, err := db.Collection("groups").InsertOne(context.Background(), tc.group)
					assert.NoError(t, err)
					tc.expected = tc.group
				}

				// テスト実行
				foundGroup, err := repo.FindGroupByGroupName(&tc.groupName)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
					assert.Nil(t, foundGroup)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, foundGroup)
					assert.Equal(t, tc.expected.GroupID, foundGroup.GroupID)
					assert.Equal(t, tc.expected.GroupName, foundGroup.GroupName)
					assert.Equal(t, tc.expected.GroupManagerID, foundGroup.GroupManagerID)
				}

				// クリーンアップ
				if tc.group != nil {
					_, err = db.Collection("groups").DeleteMany(context.Background(), bson.M{"group_name": tc.group.GroupName})
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("FindGroupByUserID", func(t *testing.T) {
		// 定数の定義
		managerID1 := entity.UserID("manager-user-id-2")
		memberID1 := entity.UserID("member1-id-2")

		testCases := []struct {
			name        string
			group       *entity.Group
			userID      entity.UserID
			expected    *entity.Group
			shouldError bool
		}{
			{
				name: "正常系: マネージャーIDでグループ検索",
				group: &entity.Group{
					GroupID:          "test-group-id-2",
					GroupName:        "TestGroup2",
					GroupPassword:    "password456",
					GroupManagerID:   managerID1,
					GroupDescription: "Test group description 2",
					GroupMembers:     []entity.UserID{memberID1, "member2-id-2"},
					GroupEvents:      []entity.EventID{"event1-id-2", "event2-id-2"},
				},
				userID:      managerID1,
				expected:    nil, // 後でセット
				shouldError: false,
			},
			{
				name: "正常系: メンバーIDでグループ検索",
				group: &entity.Group{
					GroupID:          "test-group-id-3",
					GroupName:        "TestGroup3",
					GroupPassword:    "password789",
					GroupManagerID:   "manager-user-id-3",
					GroupDescription: "Test group description 3",
					GroupMembers:     []entity.UserID{memberID1, "member2-id-3"},
					GroupEvents:      []entity.EventID{"event1-id-3", "event2-id-3"},
				},
				userID:      memberID1,
				expected:    nil, // 後でセット
				shouldError: false,
			},
			{
				name:        "異常系: 存在しないユーザーIDで検索",
				group:       nil,
				userID:      "non-existent-user-id",
				expected:    nil,
				shouldError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				if tc.group != nil {
					_, err := db.Collection("groups").InsertOne(context.Background(), tc.group)
					assert.NoError(t, err)
					tc.expected = tc.group
				}

				// テスト実行
				foundGroup, err := repo.FindGroupByUserID(tc.userID)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
					assert.Nil(t, foundGroup)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, foundGroup)
					assert.Equal(t, tc.expected.GroupID, foundGroup.GroupID)
					// 他のフィールドも必要に応じて検証
				}

				// クリーンアップ
				if tc.group != nil {
					_, err = db.Collection("groups").DeleteMany(context.Background(), bson.M{"group_id": tc.group.GroupID})
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("CreateGroup", func(t *testing.T) {
		testCases := []struct {
			name        string
			group       *entity.Group
			shouldError bool
		}{
			{
				name: "正常系: 新規グループ作成",
				group: &entity.Group{
					GroupID:          "new-group-id",
					GroupName:        "NewGroup",
					GroupPassword:    "newpassword",
					GroupManagerID:   "new-manager-id",
					GroupDescription: "New group description",
					GroupMembers:     []entity.UserID{"new-member1-id", "new-member2-id"},
					GroupEvents:      []entity.EventID{"new-event1-id", "new-event2-id"},
				},
				shouldError: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テスト実行
				createdGroup, err := repo.CreateGroup(tc.group)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, createdGroup)
					assert.Equal(t, tc.group.GroupID, createdGroup.GroupID)
					assert.Equal(t, tc.group.GroupName, createdGroup.GroupName)

					// DBに保存されていることを確認
					var savedGroup entity.Group
					err = db.Collection("groups").FindOne(context.Background(), bson.M{"group_id": tc.group.GroupID}).Decode(&savedGroup)
					assert.NoError(t, err)
					assert.Equal(t, tc.group.GroupName, savedGroup.GroupName)
					assert.Equal(t, tc.group.GroupDescription, savedGroup.GroupDescription)
				}

				// クリーンアップ
				_, err = db.Collection("groups").DeleteMany(context.Background(), bson.M{"group_id": tc.group.GroupID})
				assert.NoError(t, err)
			})
		}
	})

	t.Run("UpdateGroup", func(t *testing.T) {
		testCases := []struct {
			name         string
			initialGroup *entity.Group
			updatedGroup *entity.Group
			shouldError  bool
		}{
			{
				name: "正常系: グループ情報更新",
				initialGroup: &entity.Group{
					GroupID:          "update-group-id",
					GroupName:        "UpdateGroup",
					GroupPassword:    "updatepassword",
					GroupManagerID:   "update-manager-id",
					GroupDescription: "Update group description",
					GroupMembers:     []entity.UserID{"update-member1-id", "update-member2-id"},
					GroupEvents:      []entity.EventID{"update-event1-id", "update-event2-id"},
				},
				updatedGroup: &entity.Group{
					GroupID:          "update-group-id",
					GroupName:        "UpdateGroup",
					GroupPassword:    "updatepassword",
					GroupManagerID:   "update-manager-id",
					GroupDescription: "Updated group description",
					GroupMembers:     []entity.UserID{"update-member1-id", "update-member2-id", "update-member3-id"},
					GroupEvents:      []entity.EventID{"update-event1-id", "update-event2-id"},
				},
				shouldError: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				_, err := db.Collection("groups").InsertOne(context.Background(), tc.initialGroup)
				assert.NoError(t, err)

				// テスト実行
				updatedGroup, err := repo.UpdateGroup(tc.updatedGroup)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, updatedGroup)
					assert.Equal(t, tc.updatedGroup.GroupDescription, updatedGroup.GroupDescription)
					assert.Equal(t, len(tc.updatedGroup.GroupMembers), len(updatedGroup.GroupMembers))

					// DBから直接取得して確認
					var savedGroup entity.Group
					err = db.Collection("groups").FindOne(context.Background(), bson.M{"group_id": tc.initialGroup.GroupID}).Decode(&savedGroup)
					assert.NoError(t, err)
					assert.Equal(t, tc.updatedGroup.GroupDescription, savedGroup.GroupDescription)
					assert.Equal(t, len(tc.updatedGroup.GroupMembers), len(savedGroup.GroupMembers))
				}

				// クリーンアップ
				_, err = db.Collection("groups").DeleteMany(context.Background(), bson.M{"group_id": tc.initialGroup.GroupID})
				assert.NoError(t, err)
			})
		}
	})

	t.Run("DeleteGroup", func(t *testing.T) {
		testCases := []struct {
			name        string
			group       *entity.Group
			shouldError bool
		}{
			{
				name: "正常系: グループ削除",
				group: &entity.Group{
					GroupID:          "delete-group-id",
					GroupName:        "DeleteGroup",
					GroupPassword:    "deletepassword",
					GroupManagerID:   "delete-manager-id",
					GroupDescription: "Delete group description",
					GroupMembers:     []entity.UserID{"delete-member1-id", "delete-member2-id"},
					GroupEvents:      []entity.EventID{"delete-event1-id", "delete-event2-id"},
				},
				shouldError: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				_, err := db.Collection("groups").InsertOne(context.Background(), tc.group)
				assert.NoError(t, err)

				// テスト実行
				deletedGroup, err := repo.DeleteGroup(tc.group)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, deletedGroup)
					assert.Equal(t, tc.group.GroupID, deletedGroup.GroupID)

					// DBから削除されたことを確認
					var count int64
					count, err = db.Collection("groups").CountDocuments(context.Background(), bson.M{"group_id": tc.group.GroupID})
					assert.NoError(t, err)
					assert.Equal(t, int64(0), count)
				}
			})
		}
	})
}
