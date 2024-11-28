package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rest-app/config"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// Struktura za kreiranje korisnika
type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Nickname  string `json:"nickname"`
	Password  string `json:"password"`
}

// Struktura za ažuriranje korisnika
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// Funkcija za odgovor sa greškom
func respondWithError(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
	log.Println(message)
}

// Funkcija za kreiranje korisnika
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	// Parsiranje JSON zahteva
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Provera da li su svi podaci prisutni
	if req.FirstName == "" || req.LastName == "" || req.Nickname == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields")
		return
	}

	// Provera da li nadimak već postoji u bazi
	var exists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE nickname = $1)", req.Nickname).Scan(&exists)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	if exists {
		respondWithError(w, http.StatusConflict, "Nickname already exists")
		return
	}

	// Hešovanje lozinke
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	// Unos korisnika u bazu
	_, err = config.DB.Exec(
		"INSERT INTO users (first_name, last_name, nickname, password) VALUES ($1, $2, $3, $4)",
		req.FirstName, req.LastName, req.Nickname, string(hashedPassword),
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error inserting user into database")
		return
	}

	log.Printf("User created successfully: %s %s (%s)", req.FirstName, req.LastName, req.Nickname)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}

// Funkcija za prikaz svih korisnika sa paginacijom
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Paginacija
	limit := 10
	offset := 0

	// Preuzmi `limit` i `offset` parametre iz URL-a
	queryParams := r.URL.Query()
	if l := queryParams.Get("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}
	if o := queryParams.Get("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	// Upit za dobijanje korisnika iz baze
	rows, err := config.DB.Query("SELECT id, first_name, last_name, nickname FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	defer rows.Close()

	// Kreiranje liste korisnika
	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var firstName, lastName, nickname string
		err = rows.Scan(&id, &firstName, &lastName, &nickname)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error reading user data")
			return
		}
		users = append(users, map[string]interface{}{
			"id":        id,
			"firstName": firstName,
			"lastName":  lastName,
			"nickname":  nickname,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Funkcija za ažuriranje korisnika
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Nickname == "" && req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "No fields provided for update")
		return
	}

	// Provera da li korisnik postoji
	var exists bool
	err = config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Ažuriranje korisnika
	query := "UPDATE users SET "
	var args []interface{}
	argCounter := 1

	if req.Nickname != "" {
		query += fmt.Sprintf("nickname = $%d,", argCounter)
		args = append(args, req.Nickname)
		argCounter++
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}
		query += fmt.Sprintf("password = $%d,", argCounter)
		args = append(args, string(hashedPassword))
		argCounter++
	}

	query = query[:len(query)-1] + fmt.Sprintf(" WHERE id = $%d", argCounter)
	args = append(args, id)

	_, err = config.DB.Exec(query, args...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user")
		return
	}

	log.Printf("User with ID %s updated successfully", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

// Funkcija za brisanje korisnika
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	var userCount int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil || userCount <= 1 {
		respondWithError(w, http.StatusForbidden, "Cannot delete the only user in the database")
		return
	}

	_, err = config.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting user")
		return
	}

	log.Printf("User with ID %s deleted successfully", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
