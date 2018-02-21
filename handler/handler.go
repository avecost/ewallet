package handler

import (
	"fmt"
	"net/http"

	"github.com/avecost/ewallet/models"
)

// AppHandler is the class of Application Handler
type AppHandler struct {
	db *models.DB
}

// ErrResponse struct for Error Response (JSON)
type ErrResponse struct {
	Err     string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse struct for Success Response (JSON)
type SuccessResponse struct {
	Data interface{} `json:"data"`
}

// NewHandler create a Application Handler class
func NewHandler(db *models.DB) *AppHandler {
	return &AppHandler{db: db}
}

// Logger middleware
func (h *AppHandler) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("Request received: %v\n", req)
		next.ServeHTTP(w, req)
		fmt.Println("Request handled successfully")
	})
}
