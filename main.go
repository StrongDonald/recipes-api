// Recipes API
//
// This is a sample recipes API. You can find out more about the API at https://github.com/PacktPublishing/BuildingDistributed-Applications-in-Gin.
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 1.0.0
// Consumes:
// - application/json
// Produces:
// - application/json
// swagger:meta
package main

import (
	"context"
	"log"
	"os"
	"recipes-api/handlers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var recipesHandler *handlers.RecipesHandler

func init() {
	// Hardcoded Recipe Data
	// recipes = make([]Recipe, 0)

	// file, _ := os.ReadFile("recipes.json")

	// _ = json.Unmarshal([]byte(file), &recipes)

	// for i := 0; i < len(recipes); i++ {
	// 	recipes[i].ID = xid.New().String()
	// 	recipes[i].PublishedAt = time.Now()
	// }

	ctx := context.Background()
	client, err := mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection(("recipes"))

	log.Println("Connected to MongoDB")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	status := redisClient.Ping()
	log.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	// var listOfRecipes []interface{}

	// for _, recipe := range recipes {
	// 	listOfRecipes = append(listOfRecipes, recipe)
	// }

	// collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection(("recipes"))
	// insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}

// func SearchRecipesHandler(c *gin.Context) {
// 	tag := c.Query("tag")

// 	listOfRecipes := make([]Recipe, 0)

// 	for i := 0; i < len(recipes); i++ {
// 		found := false

// 		for _, t := range recipes[i].Tags {
// 			if strings.EqualFold(t, tag) {
// 				found = true
// 			}
// 		}

// 		if found {
// 			listOfRecipes = append(listOfRecipes,
// 				recipes[i])
// 		}
// 	}

// 	c.JSON(http.StatusOK, listOfRecipes)
// }

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Recipe not found"})

		return
	}

	recipes = append(recipes[:index], recipes[index:]...)

	c.JSON(http.StatusOK, gin.H{
		"message": "Recipe has been deleted"})
}

func main() {
	router := gin.Default()

	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	// router.GET("/recipes/search", SearchRecipesHandler)

	router.Run()
}
