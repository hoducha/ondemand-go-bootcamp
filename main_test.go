package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hoducha/ondemand-go-bootcamp/api/repositories"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testDataFile = "testdata/pokemon_data.csv"

type PokemonTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *PokemonTestSuite) SetupTest() {
	repo, err := repositories.NewPokemonRepository(testDataFile)
	if err != nil {
		log.Fatalf("Failed to initialize Pokemon repository: %v", err)
	}
	suite.router = setupRouter(repo)
}

func (suite *PokemonTestSuite) TestGetPokemonByID() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/pokemon/1", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "{\"id\":1,\"name\":\"bulbasaur\",\"image\":\"\"}", w.Body.String())
}

func (suite *PokemonTestSuite) TestGetPokemonByInvalidID() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/pokemon/invalid", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	assert.Equal(suite.T(), "{\"error\":\"Invalid ID\"}", w.Body.String())
}

func (suite *PokemonTestSuite) TestGetPokemonByIDNotFound() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/pokemon/1000", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	assert.Equal(suite.T(), "{\"error\":\"Pokemon not found\"}", w.Body.String())
}

func TestPokemonSuite(t *testing.T) {
	suite.Run(t, new(PokemonTestSuite))
}
