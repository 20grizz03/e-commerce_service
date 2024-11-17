package hendlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

type InMemoryStorage struct {
	animals []Animal
}

func New(animal []Animal) *InMemoryAnimalStorage {
	return &InMemoryAnimalStorage{animals: animal}
}

func (s *InMemoryAnimalStorage) getllAnimals(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	ch := make(chan []byte)

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			date, err := json.Marshal(s.animals)
			if err != nil {
				fmt.Println(err)
			}
			ch <- date
		}
	}()

	select {
	case <-ctx.Done():
		close(ch)
		return
	case data := <-ch:
		if data == nil {
			http.Error(w, "Unable to marshal animals", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(data)
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(s.animals)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (s *InMemoryAnimalStorage) getAnimalsByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	ch := make(chan []byte)
	idParam := chi.URLParam(r, "id")

	go func() {
		for _, animal := range s.animals {
			if fmt.Sprintf("%d", animal.ID) == idParam {
				ch <- []byte(animal.Name)
				return
			}
		}
		ch <- nil
	}()

	select {
	case <-ctx.Done():
		close(ch)
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		return
	case animal := <-ch:
		if animal == nil {
			close(ch)
			http.Error(w, "Animal not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(animal)
		if err != nil {
			http.Error(w, "Unable to marshal animal", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
	}
}

func main() {
	animals := []Animal{
		{1, "Шарик", 10},
		{2, "Маруся", 5},
	}

	storage := New(animals)

}
