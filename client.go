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

	Errors map[string]string
}

type ClientSuccessResponse struct {
	Version string  `json:"version"`
	Data    *Client `json:"data"`
}

type ClientErrorResponse struct {
	Version string  `json:"version"`
	Errors  *Client `json:"errors"`
}

type Clients struct {
	Clients []Client
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

func (c *Client) GetClientById(id int) (int, error) {
	err := c.db.QueryRow("SELECT id, name, token, address, url, reference, created_at, updated_at FROM clients WHERE id = $1", id).Scan(
		&c.ID, &c.Name, &c.Token, &c.Address, &c.URL, &c.Reference, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return 0, err
	}

	return c.ID, nil
}

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
