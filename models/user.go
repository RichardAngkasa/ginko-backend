package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID          int
	Username    string
	Password    string
	Balance     int
	BankAccount string
	CreatedAt   time.Time
}

func CreateUser(db *sql.DB, username, hashedPassword, bankAccount string) error {
	_, err := db.Exec(`
		INSERT INTO users (username, password, bank_account)
		VALUES ($1, $2, $3)
	`, username, hashedPassword, bankAccount)

	return err
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	var u User
	err := db.QueryRow(`
		SELECT id, username, password, balance, bank_account, created_at
		FROM users WHERE username=$1
	`, username).Scan(&u.ID, &u.Username, &u.Password, &u.Balance, &u.BankAccount, &u.CreatedAt)

	if err != nil {
		return nil, errors.New("user not found")
	}

	return &u, nil
}

func GetUserByID(db *sql.DB, id int) (*User, error) {
	var u User
	err := db.QueryRow(`
		SELECT id, username, password, balance, bank_account, created_at
		FROM users WHERE id=$1
	`, id).Scan(&u.ID, &u.Username, &u.Password, &u.Balance, &u.BankAccount, &u.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func UpdateBalance(db *sql.DB, id int, newBalance int) error {
	_, err := db.Exec(`UPDATE users SET balance=$1 WHERE id=$2`, newBalance, id)
	return err
}

func AddBalance(db *sql.DB, id int, amount int) error {
	_, err := db.Exec(`UPDATE users SET balance = balance + $1 WHERE id = $2`, amount, id)
	return err
}

func SubtractBalance(db *sql.DB, id int, amount int) error {
	_, err := db.Exec(`UPDATE users SET balance = balance - $1 WHERE id = $2`, amount, id)
	return err
}
