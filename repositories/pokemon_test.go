package repositories

import (
	"errors"
	"strings"
	"testing"

	"github.com/hoducha/ondemand-go-bootcamp/models"

	"github.com/stretchr/testify/assert"
)

func TestCSVRepository_GetByID(t *testing.T) {
	pokemonData := map[int]*models.Pokemon{
		1: {ID: 1, Name: "Bulbasaur"},
		2: {ID: 2, Name: "Charmander"},
	}

	repo := &CSVRepository{
		pokemonData: pokemonData,
	}

	// Test case: Existing Pokemon
	pokemon, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, pokemon)
	assert.Equal(t, 1, pokemon.ID)
	assert.Equal(t, "Bulbasaur", pokemon.Name)

	// Test case: Non-existing Pokemon
	pokemon, err = repo.GetByID(3)
	assert.Error(t, err)
	assert.Nil(t, pokemon)
	expectError := errors.New("Pokemon not found")
	assert.EqualError(t, err, expectError.Error())
}

func TestNewPokemonRepository(t *testing.T) {
	// Mock CSV data
	csvData := `1,Bulbasaur
2,Charmander`
	r := strings.NewReader(csvData)

	// Test case: Valid CSV file
	repo, err := NewPokemonRepositoryFromReader(r)
	assert.NoError(t, err)
	assert.NotNil(t, repo)

	// Test case: Invalid CSV file
	invalidCSV := "invalid csv data"
	r = strings.NewReader(invalidCSV)
	_, err = NewPokemonRepositoryFromReader(r)
	assert.Error(t, err)
}
