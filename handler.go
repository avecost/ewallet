package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// AppHandler contains a DB connection and
// handler for the Application
type AppHandler struct {
	db *DBApp
}

// NewHandler creates a new AppHandler and returns a pointer
func NewHandler(db *DBApp) *AppHandler {
	return &AppHandler{db: db}
}

// Run will initialize the routes and run the HTTP service
func (h *AppHandler) Run() {
	// check for environment our service is running on
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()
	// clients routes
	r.HandleFunc("/v1/clients", h.GetAllClients).Methods("GET")
	r.HandleFunc("/v1/clients", h.PostNewClient).Methods("POST")
	r.HandleFunc("/v1/clients/{id:[0-9]+}", h.GetClientByID).Methods("GET")
	r.HandleFunc("/v1/clients/{id:[0-9]+}", h.UpdateClientByID).Methods("PUT")
	r.HandleFunc("/v1/clients/{id:[0-9]+}", h.DeleteClientByID).Methods("DELETE")

	// wallets routes
	r.HandleFunc("/v1/{clientuuid}/wallets", h.PostNewWallet).Methods("POST")
	r.HandleFunc("/v1/{clientuuid}/wallets", h.GetAllWallet).Methods("GET")
	r.HandleFunc("/v1/{clientuuid}/wallets/{id:[0-9]+}", h.GetWalletClientUserIDHandle).Methods("GET")

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal("Server Error: ", err)
	}
}

// ValidateMiddleware will handle the security of the routes
// func (h *AppHandler) ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authorizationHeader := r.Header.Get("authorization")
// 		if authorizationHeader != "" {
// 			//TODO: get Bearer and client IP
// 			//      check if existing in our Tokens table
// 			bearerToken := strings.Split(authorizationHeader, " ")
// 			if len(bearerToken) == 2 {
// 				t := NewToken(h.db)
// 				_, err := t.ParseBearerToken(TOKEN_KEY, bearerToken)
// 				if err != nil {
// 					w.WriteHeader(http.StatusUnauthorized)
// 					var appErr []FieldError
// 					appErr = append(appErr, FieldError{Code: "400", Message: "Invalid authorization token"})
// 					json.NewEncoder(w).Encode(ErrResponse{Version: API_VERSION, Errors: appErr})

// 					return
// 				}

// 				next(w, r)
// 			}
// 		} else {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			var appErr []FieldError
// 			appErr = append(appErr, FieldError{Code: "400", Message: "An authorization header is required"})
// 			json.NewEncoder(w).Encode(ErrResponse{Version: API_VERSION, Errors: appErr})
// 			return
// 		}
// 	})
// }
