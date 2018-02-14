package main

import (
	"strings"
	"time"
)

// Client data structure
type Client struct {
	db *DBApp

	ID        int    `json:"id"`
	Name      string `json:"name"`
	Token     string `json:"token"`
	Address   string `json:"address"`
	URL       string `json:"url"`
	Reference string `json:"reference"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
	DeletedAt int    `json:"deleted_at"`
	IsActive  bool   `json:"is_active"`

	Errors map[string]string
}

// ClientUUID contains (Id, Token, IPAddr)
type ClientUUID struct {
	ID     int    `json:"id"`
	Token  string `json:"token"`
	IPAddr string `json:"ipaddr`
}

// NewClient return a pointer to a client object
func NewClient(db *DBApp) *Client {
	return &Client{db: db}
}

// Create insert a new Client to DB
func (c *Client) Create() (int, error) {
	var lastInsertID int

	cAt := time.Now().Local().Unix()
	uAt := time.Now().Local().Unix()

	err := c.db.QueryRow("INSERT INTO clients (name, token, address, url, reference, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;",
		c.Name, c.Token, c.Address, c.URL, c.Reference, cAt, uAt).Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}

	return lastInsertID, nil
}

// GetClientByID return a client record
func (c *Client) GetClientByID(id int) (int, error) {
	err := c.db.QueryRow("SELECT id, name, token, address, url, reference, created_at, updated_at, deleted_at, is_active FROM clients WHERE id = $1", id).Scan(
		&c.ID, &c.Name, &c.Token, &c.Address, &c.URL, &c.Reference, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt, &c.IsActive)
	if err != nil {
		return 0, err
	}

	return c.ID, nil
}

// UpdateClientByID return error of client update
func (c *Client) UpdateClientByID(id int) error {
	t := time.Now().Local().Unix()

	_, err := c.db.Exec("UPDATE clients SET name = $1, token = $2, address = $3, url = $4, reference = $5, updated_at = $6, is_active = $7 WHERE id = $8",
		&c.Name, &c.Token, &c.Address, &c.URL, &c.Reference, t, &c.IsActive, id)
	if err != nil {
		return err
	}

	c.GetClientByID(id)

	return nil
}

// DeleteClientByID mark record soft-delete
func (c *Client) DeleteClientByID(id int) error {
	t := time.Now().Local().Unix()

	_, err := c.db.Exec("UPDATE clients SET deleted_at = $1, is_active = $2 WHERE id = $3", t, false, id)
	if err != nil {
		return err
	}

	c.GetClientByID(id)

	return nil
}

// ValidUUID will return
func (c *Client) ValidUUID(uuid string) (*ClientUUID, bool) {
	var cUUID ClientUUID
	c.db.QueryRow("SELECT id, token, ipaddr FROM clients WHERE uuid = $1", uuid).Scan(&cUUID.ID, &cUUID.Token, &cUUID.IPAddr)
	if cUUID.ID == 0 {
		return nil, false
	}

	return &cUUID, true
}

// Validate return true or false based on validation rule
func (c *Client) Validate() bool {
	c.Errors = make(map[string]string)

	if strings.TrimSpace(c.Name) == "" {
		c.Errors["name"] = "Name is required"
	}

	if strings.TrimSpace(c.Token) == "" {
		c.Errors["token"] = "Token is required"
	}

	return len(c.Errors) == 0
}
