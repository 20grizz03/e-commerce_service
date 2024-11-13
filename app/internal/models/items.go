package models

// структура для описания самого продукта
type Item struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Price    uint64 `json:"price"`
	Quantity uint64 `json:"quantity"`
	Category string `json:"category"`
	Info     string `json:"info"`
}
