package pkg

import "time"

//Cюда мы прокидывам всек таблицы

type User struct {
	Id           int
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
