package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/router"

	_ "github.com/lib/pq"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow frontend origin
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		// Allow cookies
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Allow headers (Content-Type for JSON)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		// Allow methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Missing DB_URL environment variable")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := router.NewRouter(db)
	handler := CORSMiddleware(r)

	rows, err := db.Query("SELECT username FROM users")
	if err != nil {
		log.Fatal("Query failed:", err)
	}
	defer rows.Close()

	fmt.Println("Users in DB:")
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println("- ", name)
	}

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
