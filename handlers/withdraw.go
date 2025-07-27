package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend/middleware"
	"backend/models"
)

type WithdrawRequest struct {
	Amount int `json:"amount"`
}

func Withdraw(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		var req WithdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
			http.Error(w, "Invalid withdraw amount", http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByID(db, userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		if user.Balance < req.Amount {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		err = models.SubtractBalance(db, userID, req.Amount)
		if err != nil {
			http.Error(w, "Withdraw failed", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(`{"message":"Withdraw request submitted"}`))
	}
}
