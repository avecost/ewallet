package models

import (
	"log"
	"time"

	"github.com/rs/xid"
)

// Wallet class
type Wallet struct {
	ID        int       `json:"-"`
	Address   string    `json:"address"`
	ClientID  int       `json:"clientID"`
	UserID    int       `json:"userID"`
	Balance   float64   `json:"balance"`
	FundType  string    `json:"fundType"`
	Tag       string    `json:"tag"`
	IsActive  *bool     `json:"isActive,omitempty"`
	CreatedAT time.Time `json:"createdAt,omitempty"`
	UpdatedAT time.Time `json:"updatedAt,omitempty"`
}

// CreateWallet create wallet for a Subscribers/Users of the Client
func (db *DB) CreateWallet(wallet *Wallet) (int, error) {
	var lastInsertID int

	cAt := time.Now().Local()
	uAt := time.Now().Local()
	uid := xid.New()

	err := db.QueryRow("INSERT INTO wallets (address, client_id, user_id, balance, fund_type, tag, created_at, updated_at) "+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;",
		uid.String(), wallet.ClientID, wallet.UserID, wallet.Balance, wallet.FundType, wallet.Tag, cAt, uAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// GetWalletByID return the Wallet Object
func (db *DB) GetWalletByID(id int) (*Wallet, error) {
	var wallet Wallet
	err := db.QueryRow("SELECT address, client_id, user_id, balance, fund_type, tag FROM wallets WHERE id = $1", id).Scan(
		&wallet.Address, &wallet.ClientID, &wallet.UserID, &wallet.Balance, &wallet.FundType, &wallet.Tag)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// GetWalletByIDGUID returns a Wallet Object
func (db *DB) GetWalletByIDGUID(id int, guid string) (*Wallet, error) {
	var wallet Wallet
	err := db.QueryRow("SELECT address, client_id, user_id, balance, fund_type, tag, is_active "+
		" FROM wallets WHERE client_id = $1 AND address = $2", id, guid).Scan(
		&wallet.Address, &wallet.ClientID, &wallet.UserID, &wallet.Balance, &wallet.FundType, &wallet.Tag, &wallet.IsActive)
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

// GetAllWallet returns all Wallet of the Client
func (db *DB) GetAllWallet(id int) ([]Wallet, error) {
	rows, err := db.Query("SELECT address, client_id, user_id, balance, fund_type, tag, is_active FROM wallets WHERE client_id = $1", id)
	if err != nil {
		return nil, err
	}

	var wallets []Wallet
	for rows.Next() {
		var wallet Wallet
		err := rows.Scan(&wallet.Address, &wallet.ClientID, &wallet.UserID, &wallet.Balance, &wallet.FundType, &wallet.Tag, &wallet.IsActive)
		if err != nil {
			log.Println(err)
			continue
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

// UpdateWalletByIDGUID updates the Wallet Info
func (db *DB) UpdateWalletByIDGUID(id int, guid string, wallet *Wallet) (int64, error) {
	uAt := time.Now().Local()

	r, err := db.Exec("UPDATE wallets SET tag = $3, fund_type = $4, is_active = $5, updated_at = $6 "+
		" WHERE client_id = $1 AND address = $2;", id, guid, &wallet.Tag, &wallet.FundType, &wallet.IsActive, uAt)
	if err != nil {
		return 0, err
	}
	c, err := r.RowsAffected()
	if err != nil {
		return 0, nil
	}

	return c, nil
}

// GetWalletStatusByIDGUID check if Wallet is Active
func (db *DB) GetWalletStatusByIDGUID(id int, guid string) *bool {
	var active bool
	db.QueryRow("SELECT is_active FROM wallets WHERE client_id = $1 AND address = $2", id, guid).Scan(&active)

	return &active
}

// IsWalletActiveByIDGUID check if Wallet is active
func (db *DB) IsWalletActiveByIDGUID(id int, guid string) bool {
	var active bool
	db.QueryRow("SELECT is_active FROM wallets WHERE client_id = $1 AND address = $2", id, guid).Scan(&active)

	return active
}

// CreditWalletByIDGUID credit the wallet return newbalance or error
func (db *DB) CreditWalletByIDGUID(id int, guid string, amt float64) (*float64, *float64, error) {
	var oldBalance float64

	uAt := time.Now().Local()

	db.QueryRow("SELECT balance FROM wallets WHERE client_id = $1 AND address = $2", id, guid).Scan(&oldBalance)
	newBalance := oldBalance + amt

	_, err := db.Exec("UPDATE wallets SET balance = $3, updated_at = $4 WHERE client_id = $1 AND address = $2", id, guid, newBalance, uAt)
	if err != nil {
		return nil, nil, err
	}

	return &oldBalance, &newBalance, nil
}

// DebitWalletByIDGUID debit the wallet return newbalance or error
func (db *DB) DebitWalletByIDGUID(id int, guid string, amt float64) (*float64, *float64, error) {
	var oldBalance float64

	db.QueryRow("SELECT balance FROM wallets WHERE client_id = $1 AND address = $2", id, guid).Scan(&oldBalance)

	uAt := time.Now().Local()

	newBalance := oldBalance - amt
	_, err := db.Exec("UPDATE wallets SET balance = $3, updated_at = $4 WHERE client_id = $1 AND address = $2", id, guid, newBalance, uAt)
	if err != nil {
		return nil, nil, err
	}

	return &oldBalance, &newBalance, nil
}

// IsBalanceEnoughForDebit return bool if have enough balance for debit request
func (db *DB) IsBalanceEnoughForDebit(id int, guid string, amt float64) bool {
	var oldBalance float64

	db.QueryRow("SELECT balance FROM wallets WHERE client_id = $1 AND address = $2", id, guid).Scan(&oldBalance)
	if amt > oldBalance {
		return false
	}

	return true
}
