package models

import (
	"log"
	"time"

	"github.com/rs/xid"
)

// Client class
type Client struct {
	id        int
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	Address   string    `json:"address"`
	URL       string    `json:"url"`
	Reference string    `json:"reference"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

// CreateClient create new client in DB
func (db *DB) CreateClient(client *Client) (int, error) {
	var lastInsertID int

	cAt := time.Now().Local()
	uAt := time.Now().Local()
	uid := xid.New()

	err := db.QueryRow("INSERT INTO clients (uuid, name, token, address, url, reference, created_at, updated_at) "+
		" VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;",
		uid.String(), client.Name, client.Token, client.Address, client.URL, client.Reference, cAt, uAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// GetAllClient return all clients
func (db *DB) GetAllClient() ([]Client, error) {
	rows, err := db.Query("SELECT uuid, name, token, address, url, reference, is_active FROM clients;")
	if err != nil {
		return nil, err
	}

	var clients []Client
	for rows.Next() {
		var c Client
		err := rows.Scan(&c.UUID, &c.Name, &c.Token, &c.Address, &c.URL, &c.Reference, &c.IsActive)
		if err != nil {
			log.Println(err)
			continue
		}
		clients = append(clients, c)
	}
	return clients, nil
}

// UpdateClientByID update client info
func (db *DB) UpdateClientByID(id int, client *Client) (int64, error) {
	r, err := db.Exec("UPDATE clients SET name=$2, address=$3, url=$4, reference=$5, is_active=$6 "+
		" WHERE id=$1;", id, client.Name, client.Address, client.URL, client.Reference, client.IsActive)
	if err != nil {
		return 0, err
	}
	c, err := r.RowsAffected()
	if err != nil {
		return 0, nil
	}

	return c, nil
}

// UpdateClientByUUID update client info
func (db *DB) UpdateClientByUUID(uuid string, client *Client) (int64, error) {
	r, err := db.Exec("UPDATE clients SET name=$2, address=$3, url=$4, reference=$5, is_active=$6 "+
		" WHERE uuid=$1;", uuid, client.Name, client.Address, client.URL, client.Reference, client.IsActive)
	if err != nil {
		return 0, err
	}
	c, err := r.RowsAffected()
	if err != nil {
		return 0, nil
	}

	return c, nil
}

// GetClientByID return client object
func (db *DB) GetClientByID(id int) (*Client, error) {
	var c Client
	err := db.QueryRow("SELECT id, uuid, name, token, address, url, reference, is_active, created_at, updated_at, deleted_at "+
		" FROM clients WHERE id=$1;", id).Scan(
		&c.id, &c.UUID, &c.Name, &c.Token, &c.Address, &c.URL, &c.Reference, &c.IsActive, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetClientByUUID return client info
func (db *DB) GetClientByUUID(uuid string) (*Client, error) {
	var c Client
	err := db.QueryRow("SELECT id, uuid, name, token, address, url, reference, is_active, created_at, updated_at, deleted_at "+
		" FROM clients WHERE uuid=$1;", uuid).Scan(
		&c.id, &c.UUID, &c.Name, &c.Token, &c.Address, &c.URL, &c.Reference, &c.IsActive, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// GetClientIDByUUID return client ID
func (db *DB) GetClientIDByUUID(uuid string) (int, error) {
	var clientID int
	err := db.QueryRow("SELECT id FROM clients WHERE uuid=$1;", uuid).Scan(&clientID)
	if err != nil {
		return 0, err
	}

	return clientID, nil
}

// GetClientIDByUUIDToken if existing return the ID
func (db *DB) GetClientIDByUUIDToken(uuid, token string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM clients WHERE uuid=$1 AND token=$2;", uuid, token).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// IsValidUUIDToken verify if UUID and token is valid
func (db *DB) IsValidUUIDToken(uuid, token string) bool {
	var count int
	db.QueryRow("SELECT count(*) FROM clients WHERE uuid=$1 AND token=$2;", uuid, token).Scan(&count)
	if count == 0 {
		return false
	}

	return true
}

// IsClientActiveByUUIDToken check if client is active
func (db *DB) IsClientActiveByUUIDToken(uuid, token string) bool {
	var isActive bool
	db.QueryRow("SELECT is_active FROM clients WHERE uuid=$1 AND token=$2;", uuid, token).Scan(&isActive)

	return isActive
}
