package models

// Pokemon represents a Pokemon
type Pokemon struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Image string `json:"image"`
}
