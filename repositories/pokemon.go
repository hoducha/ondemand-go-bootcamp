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
	GetAll() []*models.Pokemon
	PersistData() error
	FilterByType(filterType string, items int, itemsPerWorker int) ([]*models.Pokemon, error)
}

// CSVRepository is a repository for getting Pokemon data from a CSV file
type CSVRepository struct {
	filePath    string
	pokemonData map[int]*models.Pokemon
}

// NewPokemonRepository returns a new PokemonRepository
func NewPokemonRepository(filePath string) (PokemonRepository, error) {
	pokemonData, err := loadData(filePath)
	if err != nil {
		return nil, err
	}

	return &CSVRepository{
		filePath:    filePath,
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

func (r *CSVRepository) FilterByType(filterType string, items int, itemsPerWorker int) ([]*models.Pokemon, error) {
	if items <= 0 {
		return nil, errors.New("items must be greater than 0")
	}
	if itemsPerWorker <= 0 {
		return nil, errors.New("itemsPerWorker must be greater than 0")
	}

	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	pokemonChan := make(chan *models.Pokemon)
	workerDone := make(chan bool)
	numOfWorkers := items / itemsPerWorker
	if items%itemsPerWorker != 0 {
		numOfWorkers++
	}

	// Start the worker pool
	for i := 0; i < numOfWorkers; i++ {
		go worker(i, reader, filterType, itemsPerWorker, pokemonChan, workerDone)
	}

	// Collect valid items from workers
	validItems := make([]*models.Pokemon, 0)
	count := 0
	hasDone := false
	for {
		select {
		case <-workerDone:
			hasDone = true
		case pokemon := <-pokemonChan:
			validItems = append(validItems, pokemon)
			count++
			if (count == items) || hasDone {
				return validItems, nil
			}
		default:
			if hasDone {
				return validItems, nil
			}
		}

	}
}

func worker(id int, reader *csv.Reader, filterType string, itemsPerWorker int, pokemonChan chan<- *models.Pokemon, workerDone chan<- bool) {
	defer func() {
		workerDone <- true
	}()

	count := 0
	for {
		if count >= itemsPerWorker {
			return
		}

		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			continue
		}

		if (filterType == "odd" && id%2 != 0) || (filterType == "even" && id%2 == 0) {
			pokemon := &models.Pokemon{
				ID:   id,
				Name: record[1],
			}
			if len(record) > 2 {
				pokemon.Image = record[2]
			}

			pokemonChan <- pokemon
			count++
		}
	}
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
func (r *CSVRepository) GetAll() []*models.Pokemon {
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
