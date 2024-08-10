package main

import (
	"fmt"
	"log"
	"main/handlers"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	fmt.Println("Welcome to building an api in GOLANG")

	// Create a router
	router := handlers.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: false,
	})

	handler := c.Handler(router)

	// listen to port
	log.Fatal(http.ListenAndServe(":5000", handler))
}
