package main

import "e-com/app/internal/handlers"

var database struct {
	Username string
	Password string
	Host     string
	Database string
}

func main() {
	handlers.StartRouter()
}
