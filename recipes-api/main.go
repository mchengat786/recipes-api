// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/Building- Distributed-Applications-in-Gin.
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Contact: Mohamed Labouardy <makkucat@gmail.com> https://labouardy.com
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
	"github.com/mchengat/recipes-api/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	ctx        context.Context
	client     *mongo.Client
	err        error
	collection *mongo.Collection
)

var (
	recipesHandler *handlers.RecipesHandler
	authHandler    *handlers.AuthHandler
)

func init() {
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("connected to mongoDB")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping(ctx)
	log.Println(status)

	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)
	authHandler = &handlers.AuthHandler{}
}

func main() {
	router := gin.Default()
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.POST("/signin", authHandler.SignInHandler)
	authorized := router.Group("/")
	authorized.Use(handlers.AuthMiddleware())
	authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
	authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipesHandler)
	authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)

	_ = router.Run(":8080")
}
