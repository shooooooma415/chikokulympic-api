package main

import (
	"context"
	"net/http"

	"chikokulympic-api/config"
	"chikokulympic-api/infrastructure/mongo/repository"
	serverV1 "chikokulympic-api/server/v1"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	_ "chikokulympic-api/docs" // Swaggerドキュメント用
)

// @title Chikokulympic API
// @version 1.0
// @description This is a Chikokulympic server API.
// @host localhost:8080
// @BasePath /
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

	// Swaggerドキュメントのエンドポイントを追加
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello Chikokulympic-api"})
	})

	// リポジトリの初期化
	userRepo := repository.NewUserRepository(db)
	groupRepo := repository.NewGroupRepository(db)

	// サーバーの初期化とルーティング設定
	authServer := serverV1.NewAuthServer(userRepo, groupRepo)
	groupServer := serverV1.NewGroupServer(groupRepo, userRepo)

	groupServer.RegisterRoutes(e)
	authServer.RegisterRoutes(e)

	port := config.GetEnvWithDefault("PORT", "8080")
	e.Start(":" + port)
}
