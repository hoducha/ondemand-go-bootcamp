package main

import (
	"log"
	"os"

	"github.com/hoducha/ondemand-go-bootcamp/api"
	"github.com/hoducha/ondemand-go-bootcamp/repositories"

	"github.com/gin-gonic/gin"
)

func setupRouter(repo repositories.PokemonRepository) *gin.Engine {
	router := gin.Default()
	api.SetupRoutes(router, repo)

	return router
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <pokemon.csv>", os.Args[0])
	}
	dataFile := os.Args[1]

	repo, err := repositories.NewPokemonRepository(dataFile)
	if err != nil {
		log.Fatalf("Failed to initialize Pokemon repository: %v", err)
	}
		
	router := setupRouter(repo)

	log.Fatal(router.Run(":8080"))
}
