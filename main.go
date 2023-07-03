package main

import (
	"log"
	"os"

	"github.com/hoducha/ondemand-go-bootcamp/api"
	"github.com/hoducha/ondemand-go-bootcamp/config"
	"github.com/hoducha/ondemand-go-bootcamp/repositories"

	"github.com/gin-gonic/gin"
)

func setupRouter(repo repositories.PokemonRepository) *gin.Engine {
	router := gin.Default()
	api.SetupRoutes(router, repo)

	return router
}

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	repo, err := repositories.NewPokemonRepository(config.Api.DataFile)
	if err != nil {
		log.Fatalf("Failed to initialize Pokemon repository: %v", err)
	}

	router := setupRouter(repo)

	log.Fatal(router.Run(":" + config.Api.Server.Port))
}
