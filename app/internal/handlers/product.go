package handlers

import (
	"context"
	"e-com/app/pkg"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
	"time"
)

// формируем фейковое хранилище
type InMemoryStorage struct {
	product []pkg.Product
}

// функция конструктор для того, чтобы сформировать объект
func New(product []pkg.Product) *InMemoryStorage {
	return &InMemoryStorage{product: product}
}

// получение всех продуктов
func (i *InMemoryStorage) getAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return
	default:
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(&i.product)
		if err != nil {
			fmt.Println(err)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

// получение продукта по ID
func (i *InMemoryStorage) getProductsById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParamStr := chi.URLParam(r, "id")

	ID, err := strconv.Atoi(idParamStr)
	if err != nil {
		http.Error(w, "Unable to convert id to int", http.StatusBadRequest)
	}

	select {
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		return
	default:
		for _, product := range i.product {
			if product.ID == uint64(ID) {
				w.Header().Set("Content-Type", "application/json")
				data, err := json.Marshal(&product)
				if err != nil {
					http.Error(w, "Unable to convert product to json", http.StatusBadRequest)
				}
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(data)
				return
			}
		}
		http.Error(w, "Product not found", http.StatusNotFound)
	}
}

func (i *InMemoryStorage) postAndPutProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return
	default:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// подумать над проверкой по обновлению ?

		var newProduct pkg.Product
		err = json.Unmarshal(body, &newProduct)
		if err != nil {
			http.Error(w, "Unable to parse body", http.StatusBadRequest)
			return
		}
		i.product = append(i.product, newProduct)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}
}

func StartRouter() {
	products := []pkg.Product{
		{1, "Шампунь", 250, 23, "Косметика", "От лысения"},
		{2, "Лосьон после бритья", 150, 54, "Косметика", "После бритья"},
	}

	storage := New(products)

	router := chi.NewRouter()

	router.Get("/product", storage.getAllProducts)
	router.Get("/products/{id}", storage.getProductsById)
	router.Post("/product", storage.postAndPutProducts)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		_ = fmt.Errorf("%v", err)
		return
	}
}
