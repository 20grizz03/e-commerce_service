package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

var database struct {
	Username string
	Password string
	Host     string
	Database string
}

func main() {
	router := chi.NewRouter()

	router.Get("/products",  {})
}
