package main

import (
	"strings"
	"time"
)

// Wallet datatype
type Wallet struct {
	db *DBApp

	ID        int     `json:"id"`
	ClientID  int     `json:"client_id"`
	UserID    int     `json:"user_id"`
	Tag       string  `json:"tag"`
	Balance   float64 `json:"balance"`
	FundType  string  `json:"fund_type"`
	CreatedAt int     `json:"created_at"`
	UpdatedAt int     `json:"updated_at"`
}

// NewWallet return a pointer to a new wallet
func NewWallet(db *DBApp) *Wallet {
	return &Wallet{db: db}
}

// Create new wallet for the user of our client
func (w *Wallet) Create() (int, error) {
	var lastInsertID int

	cAt := time.Now().Local().Unix()
	uAt := time.Now().Local().Unix()

	err := w.db.QueryRow("INSERT INTO wallets (client_id, user_id, tag, balance, fund_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;",
		w.ClientID, w.UserID, w.Tag, w.Balance, w.FundType, cAt, uAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// Validate return true or false based on validation rule
func (w *Wallet) Validate() bool {
	c.Errors = make(map[string]string)

	if strings.TrimSpace(w.ClientID .Name) == "" {
		c.Errors["name"] = "Name is required"
	}

	if strings.TrimSpace(c.Token) == "" {
		c.Errors["token"] = "Token is required"
	}

	return len(c.Errors) == 0
}
