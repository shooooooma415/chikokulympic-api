package testUtils

import (
	"context"
	"testing"

	mongoDB "chikokulympic-api/infrastructure/mongo"

	"go.mongodb.org/mongo-driver/mongo"
)

// SetupTestDB はテスト用のMongoDBデータベースを準備し、クリーンアップ関数を返します
func SetupTestDB(t *testing.T) (*mongo.Database, func()) {
	db, client, err := mongoDB.GetMongoDBConnectionWithEnvFile(".env.local")
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	cleanup := func() {
		err := db.Drop(context.Background())
		if err != nil {
			t.Logf("Warning: Failed to drop test database: %v", err)
		}

		err = mongoDB.DisconnectMongoDB(client)
		if err != nil {
			t.Logf("Warning: Failed to disconnect from MongoDB: %v", err)
		}
	}

	return db, cleanup
}
