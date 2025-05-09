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

	// FindGroupByUserIDのテストケースが削除されました

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

	t.Run("FindGroupsByUserID", func(t *testing.T) {
		// テストデータ用の定数
		commonMemberID := entity.UserID("common-member-id")
		managerID1 := entity.UserID("multi-manager-id-1")
		managerID2 := entity.UserID("multi-manager-id-2")

		// テスト用のグループデータを作成
		groups := []*entity.Group{
			{
				GroupID:          "multi-group-id-1",
				GroupName:        "MultiTestGroup1",
				GroupPassword:    "password111",
				GroupManagerID:   managerID1,
				GroupDescription: "Multi test group 1",
				GroupMembers:     []entity.UserID{commonMemberID, "member2-id-multi1"},
				GroupEvents:      []entity.EventID{"event1-id-multi1"},
			},
			{
				GroupID:          "multi-group-id-2",
				GroupName:        "MultiTestGroup2",
				GroupPassword:    "password222",
				GroupManagerID:   managerID2,
				GroupDescription: "Multi test group 2",
				GroupMembers:     []entity.UserID{commonMemberID, "member2-id-multi2"},
				GroupEvents:      []entity.EventID{"event1-id-multi2"},
			},
			{
				GroupID:          "multi-group-id-3",
				GroupName:        "MultiTestGroup3",
				GroupPassword:    "password333",
				GroupManagerID:   "another-manager-id",
				GroupDescription: "Multi test group 3",
				GroupMembers:     []entity.UserID{"different-member-id", "member2-id-multi3"},
				GroupEvents:      []entity.EventID{"event1-id-multi3"},
			},
		}

		testCases := []struct {
			name        string
			userID      entity.UserID
			expectedIDs []string
			expectedLen int
			shouldError bool
		}{
			{
				name:        "正常系: メンバーとして複数グループに所属",
				userID:      commonMemberID,
				expectedIDs: []string{"multi-group-id-1", "multi-group-id-2"},
				expectedLen: 2,
				shouldError: false,
			},
			{
				name:        "正常系: マネージャーとして所属",
				userID:      managerID1,
				expectedIDs: []string{"multi-group-id-1"},
				expectedLen: 1,
				shouldError: false,
			},
			{
				name:        "正常系: 所属グループなし",
				userID:      "non-member-id",
				expectedIDs: []string{},
				expectedLen: 0,
				shouldError: false,
			},
		}

		// テストデータをDBに挿入
		for _, group := range groups {
			_, err := db.Collection("groups").InsertOne(context.Background(), group)
			assert.NoError(t, err)
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テスト実行
				foundGroups, err := repo.FindGroupsByUserID(tc.userID)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectedLen, len(foundGroups), "Expected %d groups, got %d", tc.expectedLen, len(foundGroups))

					// 取得したグループIDが期待値に含まれているか確認
					if tc.expectedLen > 0 {
						foundIDs := make([]string, 0, len(foundGroups))
						for _, group := range foundGroups {
							foundIDs = append(foundIDs, string(group.GroupID))
						}

						for _, expectedID := range tc.expectedIDs {
							found := false
							for _, foundID := range foundIDs {
								if foundID == expectedID {
									found = true
									break
								}
							}
							assert.True(t, found, "Expected to find group ID %s but it was not returned", expectedID)
						}
					}
				}
			})
		}

		// クリーンアップ
		_, err := db.Collection("groups").DeleteMany(context.Background(), bson.M{"group_id": bson.M{"$in": []string{"multi-group-id-1", "multi-group-id-2", "multi-group-id-3"}}})
		assert.NoError(t, err)
	})

	t.Run("FindGroupByGroupID", func(t *testing.T) {
		testCases := []struct {
			name        string
			group       *entity.Group
			groupID     entity.GroupID
			expected    *entity.Group
			shouldError bool
		}{
			{
				name: "正常系: 存在するグループIDで検索",
				group: &entity.Group{
					GroupID:          "test-group-id-for-id-search",
					GroupName:        "TestGroupIDSearch",
					GroupPassword:    "password123",
					GroupManagerID:   "manager-user-id-1",
					GroupDescription: "Test group description for ID search",
					GroupMembers:     []entity.UserID{"member1-id", "member2-id"},
					GroupEvents:      []entity.EventID{"event1-id", "event2-id"},
				},
				groupID:     "test-group-id-for-id-search",
				expected:    nil, // 後でセット
				shouldError: false,
			},
			{
				name:        "異常系: 存在しないグループIDで検索",
				group:       nil,
				groupID:     "non-existent-group-id",
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
				foundGroup, err := repo.FindGroupByGroupID(tc.groupID)

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
					assert.Equal(t, tc.expected.GroupDescription, foundGroup.GroupDescription)
				}

				// クリーンアップ
				if tc.group != nil {
					_, err = db.Collection("groups").DeleteMany(context.Background(), bson.M{"group_id": tc.group.GroupID})
					assert.NoError(t, err)
				}
			})
		}
	})
}
