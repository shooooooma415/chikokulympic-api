package repository_test

import (
	"context"
	"testing"
	"time"

	"chikokulympic-api/domain/entity"
	"chikokulympic-api/infrastructure/mongo/repository"
	"chikokulympic-api/infrastructure/mongo/repository/testUtils"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestEventRepository(t *testing.T) {
	// 各テストで共通のセットアップ処理
	db, cleanup := testUtils.SetupTestDB(t)
	defer cleanup()
	repo := repository.NewEventRepository(db)

	t.Run("FindEventByEventID", func(t *testing.T) {
		// テスト用時間データ
		testTime := time.Now()

		testCases := []struct {
			name        string
			event       *entity.Event
			eventID     entity.EventID
			expected    *entity.Event
			shouldError bool
		}{
			{
				name: "正常系: 存在するイベントIDで検索",
				event: &entity.Event{
					EventID:              "test-event-id-1",
					EventTitle:           "Test Event 1",
					EventDescription:     "Test event description 1",
					EventLocationName:    "Test Location 1",
					Cost:                 1000,
					EventMessage:         "Test event message 1",
					EventAuthorID:        "test-author-id-1",
					Latitude:             35.6812,
					Longitude:            139.7671,
					EventStartDateTime:   entity.StartDateTIme(testTime),
					EventEndDateTime:     entity.EndDateTime(testTime.Add(2 * time.Hour)),
					EventClosingDateTime: entity.EventClosingDateTime(testTime.Add(-1 * time.Hour)),
					VotedMembers: []entity.VotedMember{
						{
							IsArrival:       true,
							UserID:          "test-user-id-1",
							ArrivalDateTime: testTime,
						},
					},
				},
				eventID:     "test-event-id-1",
				expected:    nil, // 後でセット
				shouldError: false,
			},
			{
				name:        "異常系: 存在しないイベントIDで検索",
				event:       nil,
				eventID:     "non-existent-event-id",
				expected:    nil,
				shouldError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				if tc.event != nil {
					_, err := db.Collection("events").InsertOne(context.Background(), tc.event)
					assert.NoError(t, err)
					tc.expected = tc.event
				}

				// テスト実行
				foundEvent, err := repo.FindEventByEventID(tc.eventID)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
					assert.Nil(t, foundEvent)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, foundEvent)
					assert.Equal(t, tc.expected.EventID, foundEvent.EventID)
					assert.Equal(t, tc.expected.EventTitle, foundEvent.EventTitle)
					assert.Equal(t, tc.expected.EventDescription, foundEvent.EventDescription)
					assert.Equal(t, tc.expected.EventAuthorID, foundEvent.EventAuthorID)
					// 必要に応じて他のフィールドも検証
				}

				// クリーンアップ
				if tc.event != nil {
					_, err = db.Collection("events").DeleteMany(context.Background(), bson.M{"_id": tc.event.EventID})
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("CreateEvent", func(t *testing.T) {
		// テスト用時間データ
		testTime := time.Now()

		testCases := []struct {
			name        string
			event       *entity.Event
			shouldError bool
		}{
			{
				name: "正常系: 新規イベント作成",
				event: &entity.Event{
					// EventIDはリポジトリで自動生成されるので設定しない
					EventTitle:           "New Event",
					EventDescription:     "New event description",
					EventLocationName:    "New Location",
					Cost:                 2000,
					EventMessage:         "New event message",
					EventAuthorID:        "new-author-id",
					Latitude:             35.6895,
					Longitude:            139.6917,
					EventStartDateTime:   entity.StartDateTIme(testTime.Add(1 * time.Hour)),
					EventEndDateTime:     entity.EndDateTime(testTime.Add(3 * time.Hour)),
					EventClosingDateTime: entity.EventClosingDateTime(testTime.Add(30 * time.Minute)),
					VotedMembers: []entity.VotedMember{
						{
							IsArrival:       false,
							UserID:          "new-user-id-1",
							ArrivalDateTime: time.Time{},
						},
					},
				},
				shouldError: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テスト実行
				createdEvent, err := repo.CreateEvent(*tc.event)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, createdEvent)

					// 自動生成されたIDが設定されていることを確認
					assert.NotEmpty(t, createdEvent.EventID)
					assert.Equal(t, tc.event.EventTitle, createdEvent.EventTitle)

					// DBに保存されていることを確認
					var savedEvent entity.Event
					err = db.Collection("events").FindOne(context.Background(), bson.M{"_id": createdEvent.EventID}).Decode(&savedEvent)
					assert.NoError(t, err)
					assert.Equal(t, tc.event.EventTitle, savedEvent.EventTitle)
					assert.Equal(t, tc.event.EventDescription, savedEvent.EventDescription)
				}

				// クリーンアップ
				_, err = db.Collection("events").DeleteMany(context.Background(), bson.M{"_id": createdEvent.EventID})
				assert.NoError(t, err)
			})
		}
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		// テスト用時間データ
		testTime := time.Now()

		testCases := []struct {
			name         string
			initialEvent *entity.Event
			updatedEvent *entity.Event
			shouldError  bool
		}{
			{
				name: "正常系: イベント情報更新",
				initialEvent: &entity.Event{
					EventID:              "update-event-id",
					EventTitle:           "Update Event",
					EventDescription:     "Update event description",
					EventLocationName:    "Update Location",
					Cost:                 3000,
					EventMessage:         "Update event message",
					EventAuthorID:        "update-author-id",
					Latitude:             35.6800,
					Longitude:            139.7700,
					EventStartDateTime:   entity.StartDateTIme(testTime),
					EventEndDateTime:     entity.EndDateTime(testTime.Add(2 * time.Hour)),
					EventClosingDateTime: entity.EventClosingDateTime(testTime.Add(-1 * time.Hour)),
					VotedMembers: []entity.VotedMember{
						{
							IsArrival:       true,
							UserID:          "update-user-id-1",
							ArrivalDateTime: testTime,
						},
					},
				},
				updatedEvent: &entity.Event{
					EventID:              "update-event-id",
					EventTitle:           "Updated Event Title",
					EventDescription:     "Updated event description",
					EventLocationName:    "Updated Location",
					Cost:                 3500,
					EventMessage:         "Updated event message",
					EventAuthorID:        "update-author-id",
					Latitude:             35.6850,
					Longitude:            139.7750,
					EventStartDateTime:   entity.StartDateTIme(testTime),
					EventEndDateTime:     entity.EndDateTime(testTime.Add(3 * time.Hour)), // 終了時間を変更
					EventClosingDateTime: entity.EventClosingDateTime(testTime.Add(-1 * time.Hour)),
					VotedMembers: []entity.VotedMember{
						{
							IsArrival:       true,
							UserID:          "update-user-id-1",
							ArrivalDateTime: testTime,
						},
						{
							IsArrival:       false,
							UserID:          "update-user-id-2",
							ArrivalDateTime: time.Time{},
						},
					},
				},
				shouldError: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				_, err := db.Collection("events").InsertOne(context.Background(), tc.initialEvent)
				assert.NoError(t, err)

				// テスト実行
				updatedEvent, err := repo.UpdateEvent(*tc.updatedEvent)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, updatedEvent)
					assert.Equal(t, tc.updatedEvent.EventTitle, updatedEvent.EventTitle)
					assert.Equal(t, tc.updatedEvent.EventDescription, updatedEvent.EventDescription)
					assert.Equal(t, tc.updatedEvent.Cost, updatedEvent.Cost)

					// DBが更新されたことを確認
					var savedEvent entity.Event
					err = db.Collection("events").FindOne(context.Background(), bson.M{"_id": tc.initialEvent.EventID}).Decode(&savedEvent)
					assert.NoError(t, err)
					assert.Equal(t, tc.updatedEvent.EventTitle, savedEvent.EventTitle)
					assert.Equal(t, len(tc.updatedEvent.VotedMembers), len(savedEvent.VotedMembers))
				}

				// クリーンアップ
				_, err = db.Collection("events").DeleteMany(context.Background(), bson.M{"_id": tc.initialEvent.EventID})
				assert.NoError(t, err)
			})
		}
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		// テスト用時間データ
		testTime := time.Now()

		testCases := []struct {
			name        string
			event       *entity.Event
			shouldError bool
		}{
			{
				name: "正常系: イベント削除",
				event: &entity.Event{
					EventID:              "delete-event-id",
					EventTitle:           "Delete Event",
					EventDescription:     "Delete event description",
					EventLocationName:    "Delete Location",
					Cost:                 4000,
					EventMessage:         "Delete event message",
					EventAuthorID:        "delete-author-id",
					Latitude:             35.6700,
					Longitude:            139.7600,
					EventStartDateTime:   entity.StartDateTIme(testTime),
					EventEndDateTime:     entity.EndDateTime(testTime.Add(2 * time.Hour)),
					EventClosingDateTime: entity.EventClosingDateTime(testTime.Add(-1 * time.Hour)),
					VotedMembers: []entity.VotedMember{
						{
							IsArrival:       true,
							UserID:          "delete-user-id-1",
							ArrivalDateTime: testTime,
						},
					},
				},
				shouldError: false,
			},
			// 必要に応じて異常系のテストケースを追加
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// テストデータのセットアップ
				_, err := db.Collection("events").InsertOne(context.Background(), tc.event)
				assert.NoError(t, err)

				// テスト実行
				deletedEvent, err := repo.DeleteEvent(*tc.event)

				// 結果の検証
				if tc.shouldError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, deletedEvent)
					assert.Equal(t, tc.event.EventID, deletedEvent.EventID)

					// DBから削除されたことを確認
					var count int64
					count, err = db.Collection("events").CountDocuments(context.Background(), bson.M{"_id": tc.event.EventID})
					assert.NoError(t, err)
					assert.Equal(t, int64(0), count)
				}
			})
		}
	})
}
