package hendlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func route() {
	router := chi.NewRouter()

	router.Get("/catalog/product", func(w http.ResponseWriter, r *http.Request) {})
}
