package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"chikokulympic-api/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoConfig struct {
	URI      string
	Database string
}

func NewMongoConfig(envFile string) *MongoConfig {
	if envFile != "" {
		config.LoadFromFileOrEnv(envFile)
	}

	uri := config.GetRequiredEnv("MONGO_URI")
	database := config.GetRequiredEnv("MONGO_DATABASE")

	return &MongoConfig{
		URI:      uri,
		Database: database,
	}
}

func ConnectMongoDB(config *MongoConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to create MongoDB client: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Printf("Connected to MongoDB at %s", config.URI)
	return client.Database(config.Database), nil
}

func DisconnectMongoDB(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect MongoDB client: %w", err)
	}

	log.Println("Disconnected from MongoDB")
	return nil
}

func GetMongoDBConnectionWithEnvFile(envFile string) (*mongo.Database, *mongo.Client, error) {
	config := NewMongoConfig(envFile)
	db, err := ConnectMongoDB(config)
	if err != nil {
		return nil, nil, err
	}

	client := db.Client()
	return db, client, nil
}
