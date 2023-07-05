package repositories_test

import (
	"fmt"
	"os"
	"testing"

	repos "github.com/hoducha/ondemand-go-bootcamp/api/repositories"

	"github.com/stretchr/testify/assert"
)

var testFilename = "../../testdata/pokemon_data_test.csv"
var numberOfPokemons = 20

// TestCSVRepository_GetByID tests the GetByID method of CSVRepository
func TestCSVRepository_GetByID(t *testing.T) {
	repo, cleanup := createTestRepository(t)
	defer cleanup()

	testcases := []struct {
		ID            int
		ExpectedName  string
		ExpectedError error
	}{
		{ID: 10, ExpectedName: "pokemon10", ExpectedError: nil},
		{ID: numberOfPokemons + 1, ExpectedName: "", ExpectedError: repos.ErrPokemonNotFound},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("ID=%d", tc.ID), func(t *testing.T) {
			pokemon, err := repo.GetByID(tc.ID)

			assert.Equal(t, tc.ExpectedError, err)
			if err == nil {
				assert.Equal(t, tc.ExpectedName, pokemon.Name)
			}
		})
	}

}

// TestCSVRepository_GetAll tests the GetAll method of CSVRepository
func TestCSVRepository_GetAll(t *testing.T) {
	repo, cleanup := createTestRepository(t)
	defer cleanup()

	pokemons := repo.GetAll()
	assert.Len(t, pokemons, numberOfPokemons)
}

// TestCSVRepository_PersistData tests the PersistData method of CSVRepository
func TestCSVRepository_PersistData(t *testing.T) {
	repo, cleanup := createTestRepository(t)
	defer cleanup()

	pokemonID := 2

	// Update a Pokemon's name
	pokemon, _ := repo.GetByID(pokemonID)
	pokemon.Name = "updated_name"

	fmt.Println(pokemon)
	fmt.Println(repo.GetAll()[pokemonID])

	err := repo.PersistData()
	assert.NoError(t, err)

	// Verify the updated data in the file
	newRepo, err := repos.NewPokemonRepository(testFilename)
	assert.NoError(t, err)

	// Verify the updated Pokemon's name
	newPokemon, _ := newRepo.GetByID(pokemonID)
	assert.Equal(t, "updated_name", newPokemon.Name)
}

func TestCSVRepository_FilterByType(t *testing.T) {
	// Create a temporary test CSV file
	filePath := "../../testdata/pokemon_data_test.csv"
	defer os.Remove(filePath)

	// Write test data to the CSV file
	err := createTestData(filePath, 100)
	assert.NoError(t, err)

	repo, err := repos.NewPokemonRepository(filePath)
	assert.NoError(t, err)
	csvRepo, ok := repo.(*repos.CSVRepository)
	assert.True(t, ok)

	testCases := []struct {
		Name           string
		FilterType     string
		Items          int
		ItemsPerWorker int
		ExectedIDs     []int
		ExpectedError  error
	}{
		{Name: "OddFilterType_Items5_ItemsPerWorker3", FilterType: "odd", Items: 5, ItemsPerWorker: 3, ExectedIDs: []int{1, 3, 5, 7, 9}, ExpectedError: nil},
		{Name: "EvenFilterType_Items4_ItemsPerWorker2", FilterType: "even", Items: 4, ItemsPerWorker: 2, ExectedIDs: []int{2, 4, 6, 8}, ExpectedError: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			pokemons, err := csvRepo.FilterByType(tc.FilterType, tc.Items, tc.ItemsPerWorker)

			assert.Equal(t, tc.ExpectedError, err)

			ids := make([]int, 0, len(pokemons))
			for _, p := range pokemons {
				ids = append(ids, p.ID)
			}
			assert.Equal(t, tc.ExectedIDs, ids)
		})
	}
}

// createTestRepository creates a PokemonRepository with test data
func createTestRepository(t *testing.T) (repos.PokemonRepository, func()) {
	err := createTestData(testFilename, numberOfPokemons)
	assert.NoError(t, err)

	repo, err := repos.NewPokemonRepository(testFilename)
	assert.NoError(t, err)

	cleanup := func() {
		os.Remove(testFilename)
	}

	return repo, cleanup
}

// createTestData creates a CSV file with test data
func createTestData(filePath string, count int) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 1; i <= count; i++ {
		line := fmt.Sprintf("%d,pokemon%d\n", i, i)
		_, err = file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}
