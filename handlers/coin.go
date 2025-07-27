package handlers

import (
	"database/sql"
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"backend/middleware"

	"github.com/lib/pq"
)

type CoinRequest struct {
	GameID string `json:"game_id"`
	Bet    int    `json:"bet"`
}

type CoinResponse struct {
	Result     string `json:"result"`
	Win        bool   `json:"win"`
	WinAmount  int    `json:"win_amount"`
	NewBalance int    `json:"new_balance"`
}

func CoinFlip(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		var req CoinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var balance int
		if err := db.QueryRow(`SELECT balance FROM users WHERE id=$1`, userID).Scan(&balance); err != nil {
			http.Error(w, "User not found", http.StatusInternalServerError)
			return
		}

		if balance < req.Bet {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		coin := "heads"
		if rand.Intn(2) == 0 {
			coin = "tails"
		}

		win := coin == "heads"
		winAmount := 0
		if win {
			winAmount = req.Bet * 2
		}

		newBalance := balance - req.Bet + winAmount

		_, err := db.Exec(`UPDATE users SET balance=$1 WHERE id=$2`, newBalance, userID)
		if err != nil {
			http.Error(w, "Failed to update balance", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			INSERT INTO spin_logs (user_id, game_id, bet, win, result, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, userID, req.GameID, req.Bet, winAmount, pq.Array([]string{coin}), time.Now())
		if err != nil {
			http.Error(w, "Failed to log spin", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(CoinResponse{
			Result:     coin,
			Win:        win,
			WinAmount:  winAmount,
			NewBalance: newBalance,
		})
	}
}
