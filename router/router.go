package router

import (
	"database/sql"
	"net/http"

	"backend/handlers"
	"backend/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Public
	r.HandleFunc("/register", handlers.Register(db)).Methods("POST")
	r.HandleFunc("/login", handlers.Login(db)).Methods("POST")

	// Protected
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.JwtAuth)

	api.HandleFunc("/fruit-spin", handlers.FruitSpin(db)).Methods("POST")
	api.HandleFunc("/coin-flip", handlers.CoinFlip(db)).Methods("POST")
	api.HandleFunc("/profile", handlers.GetProfile(db)).Methods("GET")
	api.HandleFunc("/deposit", handlers.Deposit(db)).Methods("POST")
	api.HandleFunc("/withdraw", handlers.Withdraw(db)).Methods("POST")
	api.HandleFunc("/logout", handlers.Logout()).Methods("POST")

	// Test
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running"))
	})

	return r
}
