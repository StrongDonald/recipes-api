package main

import "github.com/gin-gonic/gin"

type Recipe struct {
	Name         string   `json:"name"`
	Tags         []string `json:"tags"`
	Ingredients  []string `json:"ingredients"`
	Instructions []string `json:"instructions"`
	PublishedAt  []string `json:"publishedAt"`
}

func main() {
	router := gin.Default()

	router.Run()
}
