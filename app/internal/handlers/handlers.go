package handlers

import "net/http"

//go:generate mockery --name=Methods --output=../mocks
type Methods interface {
	getAllProducts(w http.ResponseWriter, r *http.Request)     // получение всех продуктов
	getProductsById(w http.ResponseWriter, r *http.Request)    // получение по имени
	postAndPutProducts(w http.ResponseWriter, r *http.Request) // пост запрос
}
