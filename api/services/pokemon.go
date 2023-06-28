package services

import (
	"strconv"

	"github.com/hoducha/ondemand-go-bootcamp/models"
	"github.com/hoducha/ondemand-go-bootcamp/repositories"
	"github.com/mtslzr/pokeapi-go"
)

// PokemonService is a service for Pokemon API
type PokemonService struct {
	repo repositories.PokemonRepository
}

// NewPokemonService creates a new PokemonService
func NewPokemonService(repo repositories.PokemonRepository) *PokemonService {
	return &PokemonService{repo: repo}
}

// GetByID returns a Pokemon by ID
func (s *PokemonService) GetByID(id int) (*models.Pokemon, error) {
	return s.repo.GetByID(id)
}

// UpdateImages fetches the pokemon images using the PokeAPI and updates the repository
func (s *PokemonService) UpdateImages() ([]*models.Pokemon, error) {
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
