package handlers

import (
	"net/http"
	"strconv"

	repos "github.com/hoducha/ondemand-go-bootcamp/api/repositories"
	"github.com/hoducha/ondemand-go-bootcamp/api/services"

	"github.com/gin-gonic/gin"
)

// PokemonHandler is a handler for Pokemon API
type PokemonHandler interface {
	GetByID(c *gin.Context)
	UpdateImages(c *gin.Context)
	FilterByType(c *gin.Context)
}

type pokemonHandler struct {
	service services.PokemonService
}

// NewPokemonHandler creates a new PokemonHandler
func NewPokemonHandler(repo repos.PokemonRepository) PokemonHandler {
	pokemonService := services.NewPokemonService(repo)
	return &pokemonHandler{service: pokemonService}
}

// SetupRoutes sets up routes for Pokemon API
func (h *pokemonHandler) GetByID(c *gin.Context) {
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
func (h *pokemonHandler) UpdateImages(c *gin.Context) {
	pokemons, err := h.service.UpdateImages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating images"})
		return
	}

	c.JSON(http.StatusOK, pokemons)
}

// FilterByType returns a list of Pokemon filtered by type
func (h *pokemonHandler) FilterByType(c *gin.Context) {
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
