package main

import (
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

	Err map[string]string `json:"errors"`
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

	err := w.db.QueryRow("INSERT INTO wallets (client_id, user_id, tag, fund_type, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;",
		w.ClientID, w.UserID, w.Tag, w.FundType, cAt, uAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// IsUserIDUnique return boolean if UserID is unique for the client
func (w *Wallet) IsUserIDUnique(clientID, uid int) bool {
	var cnt int
	w.db.QueryRow("SELECT count(*) FROM wallets WHERE client_id = $1 AND user_id = $2", clientID, uid).Scan(&cnt)
	if cnt != 0 {
		return false
	}
	return true
}

// IsClientWalletIDExist return bool if walletID exist on clientID
func (w *Wallet) IsClientWalletIDExist(clientID, walletID int) bool {
	var cnt int
	w.db.QueryRow("SELECT count(*) FROM wallets WHERE client_id = $1 AND user_id = $2", clientID, walletID).Scan(&cnt)
	if cnt == 0 {
		return false
	}
	return true
}

// IsAllowedToDebitClientWalletID return bool if Client Wallet ID has enough balance left
func (w *Wallet) IsAllowedToDebitClientWalletID(clientID, walletID int, amount float32) bool {
	var recID int
	var balance float32
	w.db.QueryRow("SELECT id, balance FROM wallets WHERE client_id = $1 AND user_id = $2", clientID, walletID).Scan(&recID, &balance)
	switch {
	case balance == 0:
		return false
	case amount > balance:
		return false
	default:
		return true
	}
}

// GetWalletClientUserID return the wallet info for the client user
func (w *Wallet) GetWalletClientUserID(clientID, uid int) error {
	err := w.db.QueryRow("SELECT id, client_id, user_id, tag, balance, fund_type, created_at, updated_at FROM wallets WHERE client_id = $1 AND user_id = $2", clientID, uid).Scan(
		&w.ID, &w.ClientID, &w.UserID, &w.Tag, &w.Balance, &w.FundType, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// CreditWalletBalanceClientWalletID will credit wallet balance returns previous and new balance
func (w *Wallet) CreditWalletBalanceClientWalletID(clientID, walletID int, amount float32) (float32, float32, error) {
	var recID int
	var prevBalace, newBalance float32

	err := w.db.QueryRow("SELECT id, balance FROM wallets WHERE client_id = $1 AND user_id = $2", clientID, walletID).Scan(&recID, &prevBalace)
	if err != nil {
		return 0, 0, err
	}
	newBalance = prevBalace + amount

	_, err = w.db.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", newBalance, recID)
	if err != nil {
		return 0, 0, err
	}

	return prevBalace, newBalance, nil
}

// DebitWalletBalanceClientWalletID will debit wallet balance returns previous and new balance
func (w *Wallet) DebitWalletBalanceClientWalletID(clientID, walletID int, amount float32) (float32, float32, error) {
	var recID int
	var prevBalace, newBalance float32

	err := w.db.QueryRow("SELECT id, balance FROM wallets WHERE client_id = $1 AND user_id = $2", clientID, walletID).Scan(&recID, &prevBalace)
	if err != nil {
		return 0, 0, err
	}
	newBalance = prevBalace - amount

	_, err = w.db.Exec("UPDATE wallets SET balance = $1 WHERE id = $2", newBalance, recID)
	if err != nil {
		return 0, 0, err
	}

	return prevBalace, newBalance, nil
}

// Validate return true or false based on validation rule
func (w *Wallet) Validate() bool {
	w.Err = make(map[string]string)

	// uuid => url param variable to identify client
	// token => given to client for identification need to provide for every transaction
	// user_id => the ID of the user unique for every client
	// tag => identifier of the wallet
	// balance => will always be 0.00 for every new wallet [not posted]
	// fund_type => 'default' kind of fund

	if w.UserID == 0 {
		w.Err["token"] = "UserID is required"
	}

	return len(w.Err) == 0
}
