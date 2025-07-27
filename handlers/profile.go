package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"backend/middleware"
	"backend/models"
)

type ProfileResponse struct {
	Username    string `json:"username"`
	Balance     int    `json:"balance"`
	BankAccount string `json:"bank_account"`
}

func GetProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		user, err := models.GetUserByID(db, userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(ProfileResponse{
			Username:    user.Username,
			Balance:     user.Balance,
			BankAccount: user.BankAccount,
		})
	}
}
