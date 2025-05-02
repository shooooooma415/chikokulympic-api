package main

import (
	"context"
	"net/http"

	"chikokulympic-api/config"
	"chikokulympic-api/infrastructure/mongo/repository"
	serverV1 "chikokulympic-api/server/v1"
	"chikokulympic-api/usecase"

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

	userRepo := repository.NewUserRepository(db)

	registerUserUseCase := usecase.NewRegisterUserUseCase(userRepo)
	authenticateUserUseCase := usecase.NewAuthenticateUserUseCase(userRepo)
	updateUserUseCase := usecase.NewUpdateUserUseCase(userRepo)

	authServer := serverV1.NewAuthServer(
		registerUserUseCase,
		authenticateUserUseCase,
		updateUserUseCase,
	)

	authServer.RegisterRoutes(e)

	port := config.GetEnvWithDefault("PORT", "8080")
	e.Start(":" + port)
}
