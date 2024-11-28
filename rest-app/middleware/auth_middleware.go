package middleware

import (
	"net/http"
	"rest-app/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Dohvati Authorization header
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Proveri i ukloni "Bearer " prefiks
		if len(token) < 8 || token[:7] != "Bearer " {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}
		token = token[7:]

		// Validiraj token
		claims, err := utils.ValidateJWT(token)
		if err != nil || claims == nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Nastavi ka sledećem handleru
		next.ServeHTTP(w, r)
	})
}
