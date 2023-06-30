package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Server struct{}

// (POST /storage/transaction/)
func (s Server) CreateTransaction(ctx echo.Context) error {
	data := struct {
		Message string `json:"message"`
	}{
		Message: "dummy message",
	}

	return ctx.JSON(http.StatusOK, data)
}
