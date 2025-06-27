package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"chikokulympic-api/config"
	"chikokulympic-api/infrastructure/mongo/repository"
	serverV1 "chikokulympic-api/server/v1"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "chikokulympic-api/docs"
)

// @title Chikokulympic-API
// @version 1.0
// @description This is a Chikokulympic server API.
// @host localhost:8080
// @BasePath /
var mongoConnectErr error
var db *mongo.Database

func main() {
	if os.Getenv("MONGO_URI") == "" {
		config.LoadFromFileOrEnv(".env.local")
	}

	uri := config.GetRequiredEnv("MONGO_URI")
	dbName := config.GetRequiredEnv("MONGO_DATABASE")

	log.Printf("Connecting to MongoDB: %s, Database: %s", uri, dbName)

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		mongoConnectErr = err
		log.Printf("Failed to connect to MongoDB: %v", err)
	} else {
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			mongoConnectErr = err
			log.Printf("Failed to ping MongoDB: %v", err)
		} else {
			log.Println("Successfully connected to MongoDB")
			db = client.Database(dbName)
			defer client.Disconnect(context.TODO())
		}
	}

	e := echo.New()

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello Chikokulympic-api"})
	})

	e.GET("/health", func(c echo.Context) error {
		if mongoConnectErr != nil {
			return c.JSON(http.StatusServiceUnavailable, map[string]string{"status": "unhealthy", "error": mongoConnectErr.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	if db != nil {
		userRepo := repository.NewUserRepository(db)
		groupRepo := repository.NewGroupRepository(db)
		eventRepo := repository.NewEventRepository(db)

		userServer := serverV1.NewUserServer(userRepo, groupRepo)
		groupServer := serverV1.NewGroupServer(groupRepo, userRepo)
		eventServer := serverV1.NewEventServer(eventRepo, groupRepo, userRepo)

		groupServer.RegisterRoutes(e)
		userServer.RegisterRoutes(e)
		eventServer.RegisterRoutes(e)
	}

	port := config.GetEnvWithDefault("PORT", "8080")
	log.Printf("Starting server on port %s", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
