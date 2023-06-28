package repositories

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"

	"github.com/hoducha/ondemand-go-bootcamp/models"
)

// PokemonRepository is an interface for getting Pokemon data
type PokemonRepository interface {
	GetByID(id int) (*models.Pokemon, error)
	GetAll() ([]*models.Pokemon)
	PersistData() error
}

// CSVRepository is a repository for getting Pokemon data from a CSV file
type CSVRepository struct {
	filePath string
	pokemonData map[int]*models.Pokemon
}

// NewPokemonRepository returns a new PokemonRepository
func NewPokemonRepository(filePath string) (PokemonRepository, error) {
	pokemonData, err := loadData(filePath)
	if err != nil {
		return nil, err
	}

	return &CSVRepository{
		filePath: filePath,
		pokemonData: pokemonData,
	}, nil
}

func loadData(filePath string) (map[int]*models.Pokemon, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
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
		if len(record) > 2 {
			pokemon.Image = record[2]
		}
		pokemonData[pokemon.ID] = pokemon
	}	

	return pokemonData, nil
}

// PersistData persists the Pokemon data to the CSV file
func (r *CSVRepository) PersistData() error {
	file, err := os.Create(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, pokemon := range r.pokemonData {
		err := writer.Write([]string{
			strconv.Itoa(pokemon.ID),
			pokemon.Name,
			pokemon.Image,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetAll returns all Pokemon
func (r *CSVRepository) GetAll() ([]*models.Pokemon) {
	pokemons := make([]*models.Pokemon, 0, len(r.pokemonData))
	for _, p := range r.pokemonData {
		pokemons = append(pokemons, p)
	}

	return pokemons
}

// GetByID returns a Pokemon by ID
func (r *CSVRepository) GetByID(id int) (*models.Pokemon, error) {
	pokemon, ok := r.pokemonData[id]
	if !ok {
		return nil, errors.New("Pokemon not found")
	}

	return pokemon, nil
}
