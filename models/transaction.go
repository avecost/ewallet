package models

import (
	"log"
	"time"

	"github.com/rs/xid"
)

// Transaction contains the structure of DR/CR
type Transaction struct {
	ID              int       `json:"-"`
	ClientID        int       `json:"clientId"`
	Address         string    `json:"address"`
	TransactionType string    `json:"transactionType"`
	CrAmount        float64   `json:"crAmount"`
	DrAmount        float64   `json:"drAmount"`
	OldBalance      float64   `json:"oldBalance"`
	NewBalance      float64   `json:"newBalance"`
	MethodType      string    `json:"methodType"`
	Particulars     string    `json:"particulars"`
	ReferenceCode   string    `json:"referenceCode"`
	TransactionAt   time.Time `json:"transactionAt"`
}

// CreateCreditTransaction create credit transaction for Client Subscribers/Users
func (db *DB) CreateCreditTransaction(transact *Transaction) (int, error) {
	var lastInsertID int

	uid := xid.New()
	refCode := "REF" + uid.String()

	transact.TransactionType = "cr"
	transact.ReferenceCode = refCode
	transact.TransactionAt = time.Now().Local()

	err := db.QueryRow("INSERT INTO transactions (client_id, address, transaction_type, cr_amount, "+
		" method_type, particulars, reference_code, transaction_at) "+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;",
		transact.ClientID, transact.Address, transact.TransactionType, transact.CrAmount,
		transact.MethodType, transact.Particulars, transact.ReferenceCode, transact.TransactionAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// CreateDebitTransaction create dedit transaction for Client Subscribers/Users
func (db *DB) CreateDebitTransaction(transact *Transaction) (int, error) {
	var lastInsertID int

	uid := xid.New()
	refCode := "REF" + uid.String()

	transact.TransactionType = "dr"
	transact.ReferenceCode = refCode
	transact.TransactionAt = time.Now().Local()

	err := db.QueryRow("INSERT INTO transactions (client_id, address, transaction_type, dr_amount, "+
		" method_type, particulars, reference_code, transaction_at) "+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;",
		transact.ClientID, transact.Address, transact.TransactionType, transact.DrAmount,
		transact.MethodType, transact.Particulars, transact.ReferenceCode, transact.TransactionAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// GetTransactionByID return a transaction object
func (db *DB) GetTransactionByID(id int) (*Transaction, error) {
	var t Transaction
	err := db.QueryRow("SELECT client_id, address, transaction_type, cr_amount, dr_amount, method_type, "+
		" particulars, reference_code, transaction_at "+
		" FROM transactions WHERE id = $1", id).Scan(&t.ClientID, &t.Address, &t.TransactionType, &t.CrAmount, &t.DrAmount,
		&t.MethodType, &t.Particulars, &t.ReferenceCode, &t.TransactionAt)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// GetAllTransactionByIDGUID return all Transaction for the Client e-Wallet address
func (db *DB) GetAllTransactionByIDGUID(id int, guid string) ([]Transaction, error) {
	var ts []Transaction
	rows, err := db.Query("SELECT client_id, address, transaction_type, cr_amount, dr_amount, method_type, "+
		" particulars, reference_code, transaction_at FROM transactions "+
		" WHERE client_id = $1 AND address = $2 ORDER BY transaction_at DESC", id, guid)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var t Transaction
		err := rows.Scan(&t.ClientID, &t.Address, &t.TransactionType, &t.CrAmount, &t.DrAmount,
			&t.MethodType, &t.Particulars, &t.ReferenceCode, &t.TransactionAt)
		if err != nil {
			log.Println(err)
			continue
		}
		ts = append(ts, t)
	}

	return ts, nil
}
