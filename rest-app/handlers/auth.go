package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"rest-app/config"
	"rest-app/utils"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parsiranje JSON tela zahteva
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validacija ulaznih podataka
	if req.Nickname == "" || req.Password == "" {
		log.Println("Missing required fields: nickname or password")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Dohvatanje hešovane lozinke iz baze na osnovu nadimka
	var hashedPassword string
	err = config.DB.QueryRow("SELECT password FROM users WHERE nickname = $1", req.Nickname).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found for nickname:", req.Nickname)
			http.Error(w, "Invalid nickname or password", http.StatusUnauthorized)
			return
		}
		log.Println("Database error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Provera da li se uneta lozinka poklapa sa hešovanom lozinkom
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		log.Println("Password mismatch for user:", req.Nickname)
		http.Error(w, "Invalid nickname or password", http.StatusUnauthorized)
		return
	}

	// Generisanje JWT tokena
	token, err := utils.GenerateJWT(req.Nickname)
	if err != nil {
		log.Println("Error generating JWT:", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Vraćanje tokena klijentu
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
