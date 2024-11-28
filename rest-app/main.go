package main

import (
	"log"
	"net/http"
	"rest-app/config"
	"rest-app/handlers"
	"rest-app/middleware"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Povezivanje sa bazom podataka
	config.ConnectDB()

	// Kreiranje glavnog router-a
	r := mux.NewRouter()

	// Javni endpointi
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// Zaštićeni endpointi
	api := r.PathPrefix("/users").Subrouter()
	api.Use(middleware.AuthMiddleware)
	api.HandleFunc("/", handlers.GetUsersHandler).Methods("GET")
	api.HandleFunc("/", handlers.CreateUserHandler).Methods("POST")
	api.HandleFunc("/{id}", handlers.UpdateUserHandler).Methods("PUT")
	api.HandleFunc("/{id}", handlers.DeleteUserHandler).Methods("DELETE")

	// Dodavanje CORS middleware-a
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},                   // Dozvoljeno Angular frontendu
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Dozvoljene metode
		AllowedHeaders:   []string{"Authorization", "Content-Type"},           // Dozvoljeni header-i
		AllowCredentials: true,                                                // Dozvola za kolačiće, ako su potrebni
	})

	// Primeni CORS middleware
	handler := corsMiddleware.Handler(r)

	// Pokretanje servera
	log.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", handler)
}
