package middleware

import (
	"log"
	"net/http"
	"os"

	"crypto/sha256"
	"crypto/subtle"
)

func BasicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		adminUsername := os.Getenv("ADMIN_USERNAME")
		if adminUsername == "" {
			log.Println("no admin username specified")
		}

		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			log.Println("no admin password specified")
		}

		log.Printf("adminPassword=%s\n", adminPassword)
		log.Printf("adminUsername=%s\n", adminUsername)

		username, password, ok := r.BasicAuth()
		if ok {

			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(adminUsername))
			expectedPasswordHash := sha256.Sum256([]byte(adminPassword))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			log.Printf("usernameMatch=%t\n", usernameMatch)
			log.Printf("passwordMatch=%t\n", passwordMatch)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}
