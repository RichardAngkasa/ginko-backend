package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend/middleware"
	"backend/models"
)

type DepositRequest struct {
	Amount int `json:"amount"`
}

func Deposit(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		var req DepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
			http.Error(w, "Invalid deposit amount", http.StatusBadRequest)
			return
		}

		err := models.AddBalance(db, userID, req.Amount)
		if err != nil {
			http.Error(w, "Deposit failed", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(`{"message":"Deposit successful"}`))
	}
}
