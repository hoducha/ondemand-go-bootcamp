package handlers

import (
	"net/http"
	"strconv"

	"github.com/hoducha/ondemand-go-bootcamp/repositories"

	"github.com/gin-gonic/gin"
	pokeapi "github.com/mtslzr/pokeapi-go"
)

// GetPokemonByID returns a handler function that returns a Pokemon by ID
func GetPokemonByID(repo repositories.PokemonRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		pokemon, err := repo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon not found"})
			return
		}

		c.JSON(http.StatusOK, pokemon)
	}
}

// GetPokemonDetailByID returns a handler function that returns a Pokemon detail by ID
func GetPokemonDetailByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		pokemon, err := pokeapi.Pokemon(strconv.Itoa(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon detail not found"})
			return
		}

		c.JSON(http.StatusOK, pokemon)
	}
}
