package main

import (
	"catch-all/gen/openapi"

	"github.com/labstack/echo/v4"
)

type Server struct{}

// List all pets
// (GET /pets)
func (s Server) ListPets(ctx echo.Context, params openapi.ListPetsParams) error {
	return nil
}

// Create a pet
// (POST /pets)
func (s Server) CreatePets(ctx echo.Context) error {
	return nil
}

// Info for a specific pet
// (GET /pets/{petId})
func (s Server) ShowPetById(ctx echo.Context, petId string) error {
	return nil
}
