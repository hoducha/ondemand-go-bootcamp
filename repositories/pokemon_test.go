package repositories_test

import (
	"encoding/csv"
	"io"
	"os"
	"testing"

	"github.com/hoducha/ondemand-go-bootcamp/models"
	"github.com/hoducha/ondemand-go-bootcamp/repositories"

	"github.com/stretchr/testify/assert"
)

var testDataFile = "../testdata/pokemon_data.csv"
var mockDataFile = "../testdata/pokemon_data_test.csv"

// TestCSVRepository_GetByID tests the GetByID method of CSVRepository
func TestCSVRepository_GetByID(t *testing.T) {
	repo, cleanup := createTestRepository(t)
	defer cleanup()

	// Test case: Existing Pokemon
	pokemon, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, pokemon)
	assert.Equal(t, 1, pokemon.ID)
	assert.Equal(t, "bulbasaur", pokemon.Name)

	// Test case: Non-existing Pokemon
	pokemon, err = repo.GetByID(1000)
	assert.Error(t, err)
	assert.Nil(t, pokemon)
	assert.EqualError(t, err, "Pokemon not found")
}

// TestCSVRepository_GetAll tests the GetAll method of CSVRepository
func TestCSVRepository_GetAll(t *testing.T) {
	repo, cleanup := createTestRepository(t)
	defer cleanup()

	pokemons := repo.GetAll()
	assert.Len(t, pokemons, 3)

	expectedPokemons := []*models.Pokemon{
		{ID: 1, Name: "bulbasaur"},
		{ID: 2, Name: "ivysaur"},
		{ID: 3, Name: "venusaur"},
	}

	assert.ElementsMatch(t, expectedPokemons, pokemons)
}

// TestCSVRepository_PersistData tests the PersistData method of CSVRepository
func TestCSVRepository_PersistData(t *testing.T) {
	repo, cleanup := createTestRepository(t)
	defer cleanup()

	// Update a Pokemon's name
	pokemon, _ := repo.GetByID(2)
	pokemon.Name = "updated_name"

	err := repo.PersistData()
	assert.NoError(t, err)

	// Verify the updated data in the file
	file, err := os.Open(mockDataFile)
	assert.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	assert.NoError(t, err)

	// Verify the updated Pokemon's name
	assert.Equal(t, "updated_name", records[1][1])
}

// createTestRepository creates a CSVRepository for testing
func createTestRepository(t *testing.T) (repositories.PokemonRepository, func()) {
	// Copy the mock data file to a temporary location
	tmpFile, err := copyFile(testDataFile, mockDataFile)
	assert.NoError(t, err)

	repo, err := repositories.NewPokemonRepository(tmpFile)
	assert.NoError(t, err)

	cleanup := func() {
		os.Remove(tmpFile)
	}

	return repo, cleanup
}

// copyFile copies a file to a destination path and returns the destination path
func copyFile(src, dst string) (string, error) {
	in, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return "", err
	}

	return dst, nil
}
