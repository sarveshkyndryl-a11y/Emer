package main

import (
	
	"log"
	"net/http"

	"github.com/joho/godotenv"
)
func main() {
if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	
	db.Connect()
	err := http.ListenAndServe(":8080", middlewares.CORS(router.Router()))
	if err != nil {
		panic(err)
	}
}