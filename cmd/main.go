package main

import (
	"context"
	"net/http"

	"chikokulympic-api/config"
	"chikokulympic-api/infrastructure/mongo/repository"
	serverV1 "chikokulympic-api/server/v1"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config.LoadEnvFileOrDefault(".env.local")

	uri := config.GetRequiredEnv("MONGO_URI")
	dbName := config.GetRequiredEnv("MONGO_DATABASE")

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	db := client.Database(dbName)

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello Chikokulympic-api"})
	})

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)

	// サーバーの初期化とルーティング設定
	authServer := serverV1.NewAuthServer(userRepo)
	authServer.RegisterRoutes(e)

	port := config.GetEnvWithDefault("PORT", "8080")
	e.Start(":" + port)
}
