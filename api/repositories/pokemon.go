package repositories

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"

	"github.com/hoducha/ondemand-go-bootcamp/api/models"
)

// Define errors
var (
	ErrPokemonNotFound = errors.New("Pokemon not found")
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

// FilterByType returns a list of Pokemon filtered by type (odd or even)
// The function uses multiple workers to read the CSV file concurrently
func (r *CSVRepository) FilterByType(filterType string, items int, itemsPerWorker int) ([]*models.Pokemon, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	numOfWorkers := items / itemsPerWorker
	if items%itemsPerWorker != 0 {
		numOfWorkers++
	}

	fileSize := fileInfo.Size()
	chunkSize := fileSize / int64(numOfWorkers)

	pokemonChan := make(chan *models.Pokemon, items)
	workerDone := make(chan bool)

	// Start the worker pool
	for i := 0; i < numOfWorkers; i++ {
		workerReader := csv.NewReader(file)
		startOffset := int64(i) * chunkSize
		endOffset := startOffset + chunkSize

		// Skip lines before the start offset
		if startOffset > 0 {
			_, err := file.Seek(startOffset, io.SeekStart)
			if err != nil {
				return nil, err
			}
			line, err := workerReader.Read()
			if err != nil && err != io.EOF {
				return nil, err
			}
			// Adjust start offset to the beginning of the line
			if len(line) > 0 {
				startOffset = startOffset + int64(len(line[0]))
			}
		}

		// Adjust end offset if it falls in the middle of a line
		_, err := file.Seek(endOffset, io.SeekStart)
		if err != nil {
			return nil, err
		}
		line, _ := workerReader.Read()
		if len(line) > 0 {
			endOffset += int64(len(line[0]))
		}

		// Create a separate CSV reader for each worker
		workerFile := io.NewSectionReader(file, startOffset, endOffset-startOffset)
		workerReader = csv.NewReader(workerFile)

		go worker(workerReader, filterType, pokemonChan, workerDone)
	}

	// Collect valid items from workers
	validItems := make([]*models.Pokemon, 0)
	count := 0
	hasDone := false
	for count < items {
		select {
		case <-workerDone:
			numOfWorkers--
			if numOfWorkers == 0 {
				hasDone = true
			}
		case pokemon := <-pokemonChan:
			validItems = append(validItems, pokemon)
			count++
		default:
			if hasDone {
				return validItems, nil
			}
		}
	}

	return validItems, nil
}

func worker(reader *csv.Reader, filterType string, pokemonChan chan<- *models.Pokemon, workerDone chan<- bool) {
	defer func() {
		workerDone <- true
	}()

	count := 0
	for {
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
		return nil, ErrPokemonNotFound
	}

	return pokemon, nil
}
