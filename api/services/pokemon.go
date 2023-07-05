package services

import (
	"strconv"

	"github.com/hoducha/ondemand-go-bootcamp/api/models"
	repos "github.com/hoducha/ondemand-go-bootcamp/api/repositories"

	"github.com/mtslzr/pokeapi-go"
)

// PokemonService is a service for Pokemon API
type PokemonService interface {
	GetByID(id int) (*models.Pokemon, error)
	UpdateImages() ([]*models.Pokemon, error)
	FilterByType(filterType string, items int, itemsPerWorker int) ([]*models.Pokemon, error)
}

type pokemonService struct {
	repo repos.PokemonRepository
}

// NewPokemonService creates a new PokemonService
func NewPokemonService(repo repos.PokemonRepository) PokemonService {
	return &pokemonService{repo: repo}
}

// GetByID returns a Pokemon by ID
func (s *pokemonService) GetByID(id int) (*models.Pokemon, error) {
	return s.repo.GetByID(id)
}

// UpdateImages fetches the pokemon images using the PokeAPI and updates the repository
func (s *pokemonService) UpdateImages() ([]*models.Pokemon, error) {
	pokemons := s.repo.GetAll()
	for _, pokemon := range pokemons {
		data, err := pokeapi.Pokemon(strconv.Itoa(pokemon.ID))
		if err != nil {
			return nil, err
		}
		pokemon.Image = data.Sprites.FrontDefault
	}

	s.repo.PersistData()

	return pokemons, nil
}

// FilterByType returns a list of Pokemon filtered by type
func (s *pokemonService) FilterByType(filterType string, items int, itemsPerWorker int) ([]*models.Pokemon, error) {
	return s.repo.FilterByType(filterType, items, itemsPerWorker)
}
