package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Zameni vrednost sa lozinkom iz baze
	hashedPassword := "$2a$10$CwTycUXWue0Thq9StjUM0uJ8ZhSxtAfC6cmFZQTT9DHcyUJGib8Iq"
	password := "password123"

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println("Password mismatch:", err)
	} else {
		fmt.Println("Password matched!")
	}
}
