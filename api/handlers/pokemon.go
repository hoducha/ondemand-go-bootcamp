package handlers

import (
	"net/http"
	"strconv"

	"github.com/hoducha/ondemand-go-bootcamp/api/services"
	"github.com/hoducha/ondemand-go-bootcamp/repositories"

	"github.com/gin-gonic/gin"
)

// PokemonHandler is a handler for Pokemon API
type PokemonHandler struct {
	service *services.PokemonService
}

// NewPokemonHandler creates a new PokemonHandler
func NewPokemonHandler(repo repositories.PokemonRepository) *PokemonHandler {
	pokemonService := services.NewPokemonService(repo)
	return &PokemonHandler{service: pokemonService}
}

// SetupRoutes sets up routes for Pokemon API
func (h *PokemonHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	pokemon, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon not found"})
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

// UpdateImages fetches the pokemon images using the PokeAPI and updates the repository
func (h *PokemonHandler) UpdateImages(c *gin.Context) {
	pokemons, err := h.service.UpdateImages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating images"})
		return
	}

	c.JSON(http.StatusOK, pokemons)
}

// FilterByType returns a list of Pokemon filtered by type
func (h *PokemonHandler) FilterByType(c *gin.Context) {
	filterType := c.Query("type")
	items, err := strconv.Atoi(c.Query("items"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid items"})
		return
	}
	itemsPerWorker, err := strconv.Atoi(c.Query("itemsPerWorker"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid itemsPerWorker"})
		return
	}

	pokemons, err := h.service.FilterByType(filterType, items, itemsPerWorker)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error filtering by type"})
		return
	}

	c.JSON(http.StatusOK, pokemons)
}
