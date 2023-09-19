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
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)

	file, _ := os.ReadFile("recipes.json")

	_ = json.Unmarshal([]byte(file), &recipes)

	for i := 0; i < len(recipes); i++ {
		recipes[i].ID = xid.New().String()
		recipes[i].PublishedAt = time.Now()
	}

}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})

		return
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)

	c.JSON(http.StatusOK, recipe)
}

// swagger:operation GET /recipes recipes listRecipes
// Returns list of recipes
// ---
// produces:
// - application/json
// responses:
// '200':
// description: Successful operation
func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	var recipe Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})

		return
	}

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

	id = recipes[index].ID

	recipes[index] = recipe

	recipes[index].ID = id
	recipes[index].PublishedAt = time.Now()

	c.JSON(http.StatusOK, recipes)
}

func main() {
	router := gin.Default()

	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)

	router.Run()
}
