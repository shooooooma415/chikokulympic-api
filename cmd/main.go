package main

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := os.Getenv("MONGO_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	db := client.Database("chikokulympic_prod")
	collection := db.Collection("hoge")

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		var result bson.M
		err := collection.FindOne(context.TODO(), bson.M{}).Decode(&result)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, result)
	})

	e.Start(":8080")
}
