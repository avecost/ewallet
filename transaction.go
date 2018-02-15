package main

import "time"

// Transaction contains the structure of DR/CR
type Transaction struct {
	db *DBApp

	ID              int     `json:"transaction_id"`
	ClientID        int     `json:"client_id"`
	WalletID        int     `json:"user_id"`
	TransactionType string  `json:"transaction_type"`
	CrAmount        float32 `json:"credit_amount"`
	DrAmount        float32 `json:"debit_amount"`
	OldBalance      float32 `json:"old_balance"`
	NewBalance      float32 `json:"new_balance"`
	MethodType      string  `json:"method_type"`
	Particulars     string  `json:"particulars"`
	CreatedAt       int     `json:"created_at"`
	UpdatedAt       int     `json:"updated_at"`
}

// NewTransaction return a Transaction Object
func NewTransaction(db *DBApp) *Transaction {
	return &Transaction{db: db}
}

// Credit will insert a credit transaction for Client Wallet ID
func (t *Transaction) Credit() (int, error) {
	var lastInsertID int

	cAt := time.Now().Local().Unix()
	uAt := time.Now().Local().Unix()

	err := t.db.QueryRow("INSERT INTO transactions (client_id, user_id, transaction_type, cr_amount, method_type, particulars, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;",
		&t.ClientID, &t.WalletID, &t.TransactionType, &t.CrAmount, &t.MethodType, &t.Particulars, cAt, uAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

// Debit will insert a debit transaction for Client Wallet ID
func (t *Transaction) Debit() (int, error) {
	var lastInsertID int

	cAt := time.Now().Local().Unix()
	uAt := time.Now().Local().Unix()

	err := t.db.QueryRow("INSERT INTO transactions (client_id, user_id, transaction_type, dr_amount, method_type, particulars, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;",
		&t.ClientID, &t.WalletID, &t.TransactionType, &t.DrAmount, &t.MethodType, &t.Particulars, cAt, uAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}
