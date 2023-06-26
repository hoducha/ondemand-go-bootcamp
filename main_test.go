package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testDataFile = "pokemon_data.csv"

type PokemonTestSuite struct {
	suite.Suite
	router *gin.Engine
}

func (suite *PokemonTestSuite) SetupTest() {
	suite.router = setupRouter(testDataFile)
}

func (suite *PokemonTestSuite) TestGetPokemonByID() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/pokemon/1", nil)
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
	assert.Equal(suite.T(), "{\"id\":1,\"name\":\"bulbasaur\"}", w.Body.String())
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
