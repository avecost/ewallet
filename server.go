package ewallet

import (
	"fmt"
	"log"
	"net/http"

	"github.com/avecost/ewallet/handler"
	"github.com/avecost/ewallet/models"

	"github.com/gorilla/mux"
)

// Server is the Application Class
type Server struct {
	db *models.DB
}

// NewServer create our server
func NewServer(conn string) *Server {
	c, err := models.NewDB(conn)
	if err != nil {
		panic(err)
	}

	return &Server{db: c}
}

// Run the main loop of the server
func (s *Server) Run(addr string) {
	// load the routes
	s.router(addr)
	// make sure we close the db session
	defer s.db.Close()
}

func (s *Server) router(port string) {
	certPath := "server.pem"
	keyPath := "server.key"

	// create handler object
	h := handler.NewHandler(s.db)

	r := mux.NewRouter()

	// client routes
	r.Handle("/v1/clients", h.Logger(http.HandlerFunc(h.ClientGetAllHandler))).Methods("GET")
	r.HandleFunc("/v1/clients", h.ClientPostHandler).Methods("POST")
	r.HandleFunc("/v1/clients/{uuid}", h.ClientGetHandler).Methods("GET")
	r.HandleFunc("/v1/clients/{uuid}", h.ClientPutHandler).Methods("PUT")

	// wallet routes
	r.Handle("/v1/{uuid}/wallets", h.WithTokenMiddleware(http.HandlerFunc(h.WalletGetAllHandler))).Methods("GET")
	r.Handle("/v1/{uuid}/wallets", h.WithTokenMiddleware(http.HandlerFunc(h.WalletPostHandler))).Methods("POST")
	r.Handle("/v1/{uuid}/wallets/{guid}", h.WithTokenMiddleware(http.HandlerFunc(h.WalletGetHandler))).Methods("GET")
	r.Handle("/v1/{uuid}/wallets/{guid}", h.WithTokenMiddleware(http.HandlerFunc(h.WalletPutHandler))).Methods("PUT")

	// transaction routes
	r.Handle("/v1/{uuid}/transaction/{guid}", h.WithTokenMiddleware(http.HandlerFunc(h.GetAllTransactionHandler))).Methods("GET")
	r.Handle("/v1/{uuid}/transaction/{guid}/credit", h.WithTokenMiddleware(http.HandlerFunc(h.CreditPostHandler))).Methods("POST")
	r.Handle("/v1/{uuid}/transaction/{guid}/debit", h.WithTokenMiddleware(http.HandlerFunc(h.DebitPostHandler))).Methods("POST")

	// inform that we are live
	fmt.Println("e-Wallet is running on port: ", port)

	// err := http.ListenAndServe(port, r)
	err := http.ListenAndServeTLS(port, certPath, keyPath, r)
	if err != nil {
		log.Fatal("Server Error: ", err)
	}
}
