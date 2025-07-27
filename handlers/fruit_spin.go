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

type FruitSpinRequest struct {
	GameID string `json:"game_id"`
	Bet    int    `json:"bet"`
}

type FruitSpinResponse struct {
	Result     []string `json:"result"`
	Win        int      `json:"win"`
	NewBalance int      `json:"new_balance"`
	IsJackpot  bool     `json:"is_jackpot"`
}

var fruits = []string{"cherry", "banana", "lemon", "apple"}

func FruitSpin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)

		var req FruitSpinRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		var balance int
		err := db.QueryRow(`SELECT balance FROM users WHERE id=$1`, userID).Scan(&balance)
		if err != nil || balance < req.Bet {
			http.Error(w, "Insufficient balance", http.StatusBadRequest)
			return
		}

		result := []string{
			fruits[rand.Intn(len(fruits))],
			fruits[rand.Intn(len(fruits))],
			fruits[rand.Intn(len(fruits))],
		}

		win := 0
		isJackpot := result[0] == result[1] && result[1] == result[2]
		if isJackpot {
			win = req.Bet * 10
		}

		newBalance := balance - req.Bet + win
		_, err = db.Exec(`UPDATE users SET balance=$1 WHERE id=$2`, newBalance, userID)
		if err != nil {
			http.Error(w, "Failed to update balance", http.StatusInternalServerError)
			return
		}

		_, err = db.Exec(`
			INSERT INTO spin_logs (user_id, game_id, bet, win, result, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, userID, req.GameID, req.Bet, win, pq.Array(result), time.Now())
		if err != nil {
			http.Error(w, "Failed to log spin", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(FruitSpinResponse{
			Result:     result,
			Win:        win,
			NewBalance: newBalance,
			IsJackpot:  isJackpot,
		})
	}
}
