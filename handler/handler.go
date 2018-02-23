package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

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

// WithTokenMiddleware requires the request to have a valid token before allowing to proceed
func (h *AppHandler) WithTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		// Get token from the Authorization header
		// format: Authorization: Bearer
		tokens, ok := r.Header["Authorization"]
		if ok && len(tokens) >= 1 {
			token = tokens[0]
			token = strings.TrimPrefix(token, "Bearer ")
		}

		// If the token is empty...
		if token == "" {
			// If we get here, the required token is missing
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		uuid := vars["uuid"]
		// check if uuid and token do exist
		if !h.db.IsValidUUIDToken(uuid, token) {
			// If we get here, the required token is missing
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// check if client is Active
		if !h.db.IsClientActiveByUUIDToken(uuid, token) {
			// If we get here, the required token is missing
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

		// // Now parse the token
		// parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		//     // Don't forget to validate the alg is what you expect:
		//     if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//         msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		//         return nil, msg
		//     }
		//     return a.encryptionKey, nil
		// })
		// if err != nil {
		//     http.Error(w, "Error parsing token", http.StatusUnauthorized)
		//     return
		// }

		// // Check token is valid
		// if parsedToken != nil && parsedToken.Valid {
		//     // Everything worked! Set the user in the context.
		//     context.Set(r, "user", parsedToken)
		//     next.ServeHTTP(w, r)
		//     fmt.Println("test")
		// }

		// // Token is invalid
		// http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		// return
	})
}
