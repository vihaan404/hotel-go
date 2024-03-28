package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/vihaan404/hotel-go/api"
	"github.com/vihaan404/hotel-go/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dburi    = "mongodb+srv://vihaan2005yadav:fU1LyD9y2QJKZEDF@cluster0.mxzwc.mongodb.net/"
	dbName   = "hotel-reservation"
	userColl = "users"
)

// Create a new fiber instance with custom config
var config = fiber.Config{
	// Override default error handler
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

// ...
func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}
	// handler intilization
	dbStore := db.NewMongoUserStore(client)
	userHandler := api.NewUserHandler(dbStore)

	listenAddr := flag.String("listenAddr", ":5000", "specify the port")
	flag.Parse()
	app := fiber.New(config)

	appv1 := app.Group("api/v1")
	appv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	appv1.Post("/user", userHandler.HandlePostUser)
	appv1.Get("/user", userHandler.HandleGetUsers)

	appv1.Get("/user/:id", userHandler.HandleGetUser)
	app.Listen(*listenAddr)
}
