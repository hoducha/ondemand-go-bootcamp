package repositories

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/hoducha/ondemand-go-bootcamp/models"
)

// PokemonRepository is an interface for getting Pokemon data
type PokemonRepository interface {
	GetByID(id int) (*models.Pokemon, error)
}

// CSVRepository is a repository for getting Pokemon data from a CSV file
type CSVRepository struct {
	pokemonData map[int]*models.Pokemon
}

// NewPokemonRepository returns a new PokemonRepository
func NewPokemonRepository(filePath string) (PokemonRepository, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewPokemonRepositoryFromReader(file)
}

// NewPokemonRepositoryFromReader returns a new PokemonRepository from a reader
func NewPokemonRepositoryFromReader(reader io.Reader) (PokemonRepository, error) {
	records, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, err
	}

	pokemonData := make(map[int]*models.Pokemon)
	for _, record := range records {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}
		pokemon := &models.Pokemon{
			ID:   id,
			Name: record[1],
		}
		pokemonData[pokemon.ID] = pokemon
	}

	return &CSVRepository{
		pokemonData: pokemonData,
	}, nil
}

// GetByID returns a Pokemon by ID
func (r *CSVRepository) GetByID(id int) (*models.Pokemon, error) {
	pokemon, ok := r.pokemonData[id]
	if !ok {
		return nil, errors.New("Pokemon not found")
	}

	return pokemon, nil
}
