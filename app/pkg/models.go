package pkg

import "time"

// Cюда мы прокидывам все таблицы
type User struct {
	Id           int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// структура Product
type Product struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Price    uint32 `json:"price"`
	Quantity uint64 `json:"quantity"`
	Category string `json:"category"`
	Info     string `json:"info"`
}
