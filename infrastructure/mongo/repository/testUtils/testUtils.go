package testUtils

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	mongoDB "chikokulympic-api/infrastructure/mongo"

	"go.mongodb.org/mongo-driver/mongo"
)

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", os.ErrNotExist
		}
		dir = parentDir
	}
}

func SetupTestDB(t *testing.T) (*mongo.Database, func()) {
	rootDir, err := findProjectRoot()
	if err != nil {
		t.Fatalf("Failed to find project root: %v", err)
	}

	envPath := filepath.Join(rootDir, ".env.local")

	db, client, err := mongoDB.GetMongoDBConnectionWithEnvFile(envPath)
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
