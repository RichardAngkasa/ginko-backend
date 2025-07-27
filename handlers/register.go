package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend/models"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	BankAccount string `json:"bank_account"`
}

func Register(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		err = models.CreateUser(db, req.Username, string(hashedPassword), req.BankAccount)
		if err != nil {
			http.Error(w, "User already exists or DB error", http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"user created"}`))
	}
}
