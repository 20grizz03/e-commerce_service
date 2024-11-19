package handlers

import (
	"context"
	"e-com/app/pkg"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"time"
)

type InMemoryStorage struct {
	product []pkg.Product
}

func New(product []pkg.Product) *InMemoryStorage {
	return &InMemoryStorage{product: product}
}

func (i *InMemoryStorage) getAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	ch := make(chan []byte)

	go func() {
		select {
		case <-ctx.Done():
			return
		default:
			date, err := json.Marshal(&i.product)
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
	data, err := json.Marshal(&i.product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(data)
}

func (i *InMemoryStorage) getProductsByName(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	ch := make(chan []byte)
	idParamStr := chi.URLParam(r, "id")

	ID, err := strconv.Atoi(idParamStr)
	if err != nil {
		http.Error(w, "Unable to convert id to int", http.StatusBadRequest)
	}

	go func() {
		for _, product := range i.product {
			if product.ID == uint64(ID) {
				ch <- []byte(product.Name)
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
	case product := <-ch:
		if product == nil {
			close(ch)
			http.Error(w, "Animal not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(string(product))
		if err != nil {
			http.Error(w, "Unable to marshal animal", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(data)
	}
}

// удаление продукта по ID
func (i *InMemoryStorage) deleteProductById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParamStr := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(idParamStr)
	if err != nil {
		http.Error(w, "Unable to convert id to int", http.StatusBadRequest)
	}

	ch := make(chan bool)
	go func() {
		// решение только для среза, но надо переделать чтобы удалялось из БД
		deleted := false
		var filteredProducts []pkg.Product // создаем новый срез продуктов

		for _, product := range i.product {
			if product.ID == uint64(ID) {
				deleted = true
				continue // Пропускаем элемент, который нужно удалить
			}
			filteredProducts = append(filteredProducts, product)
		}
		// обновляем оригинальный срез только если элемент был найден и удален
		if deleted {
			i.product = filteredProducts
			ch <- true
		} else {
			ch <- false
		}
	}()

	select {
	case <-ctx.Done():
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		return
	case deleted := <-ch:
		if !deleted {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Product deleted successfully"))
	}
}

// обновление продукта по ID
func (i *InMemoryStorage) updateProductById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	idParamStr := chi.URLParam(r, "id")
	ID, err := strconv.Atoi(idParamStr)
	if err != nil {
		http.Error(w, "Unable to convert id to int", http.StatusBadRequest)
	}

	// получаем данные из json файла и
	var updatedProduct pkg.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	updatedProduct.ID = uint64(ID)

	ch := make(chan bool)
	go func() {
		// Ищем и обновляем продукт
		for index, product := range i.product {
			if product.ID == uint64(ID) {
				i.product[index] = updatedProduct
				ch <- true // Продукт найден и обновлен
				return
			}
		}
		ch <- false // Продукт с данным ID не найден
	}()

	select {
	case <-ctx.Done():
		close(ch)
		http.Error(w, "Request timed out", http.StatusGatewayTimeout)
		return
	case updated := <-ch:
		if !updated {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Product updated successfully"))
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
	router.Get("/products/{id}", storage.getProductsByName)
	router.Delete("/products/{id}", storage.deleteProductById)
	router.Put("/products/{id}", storage.updateProductById)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		_ = fmt.Errorf("%v", err)
		return
	}
}
