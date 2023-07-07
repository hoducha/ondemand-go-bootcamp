package api

import (
	"github.com/hoducha/ondemand-go-bootcamp/api/handlers"
	repos "github.com/hoducha/ondemand-go-bootcamp/api/repositories"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the routes for the API
func SetupRoutes(router *gin.Engine, repo repos.PokemonRepository) {
	pokemonHandler := handlers.NewPokemonHandler(repo)

	v1 := router.Group("/v1")
	{
		v1.GET("/pokemon/update-images", pokemonHandler.UpdateImages)
		v1.GET("/pokemon/filter", pokemonHandler.FilterByType)
		v1.GET("/pokemon/:id", pokemonHandler.GetByID)
	}
}
