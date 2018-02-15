package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/avecost/ewallet/response"
)

// GetAllWallet returns all wallets for the client
func (h *AppHandler) GetAllWallet(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	clientID, Ok := v["clientid"]
	if !Ok {
		response.BadRequest(w, "Invalid request")
		return
	}
	//TODO: check if client id is existing

	fmt.Println(clientID)
}

// GetWalletClientUserIDHandle handle the request to return Wallet of Client UserID
func (h *AppHandler) GetWalletClientUserIDHandle(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	clientUUID := v["clientuuid"]
	userID := v["id"]

	u := NewClient(h.db)
	uuid, Ok := u.ValidUUID(clientUUID)
	if !Ok {
		response.BadRequest(w, "Invalid Client Key")
		return
	}

	if userID == "" {
		response.BadRequest(w, "UserID is required")
		return
	}

	uid, err := strconv.Atoi(userID)
	if err != nil {
		response.BadRequest(w, "UserID is required")
		return
	}

	q := r.URL.Query()
	token := q.Get("token")
	if token != uuid.Token || token == "" {
		response.BadRequest(w, "Invalid token")
		return
	}

	wallet := NewWallet(h.db)
	wallet.GetWalletClientUserID(uuid.ID, uid)

	response.JSON(w, wallet, http.StatusOK)
}

// PostNewWallet create a new Wallet
func (h *AppHandler) PostNewWallet(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	clientUUID := v["clientuuid"]

	u := NewClient(h.db)
	uuid, Ok := u.ValidUUID(clientUUID)
	if !Ok {
		response.BadRequest(w, "Invalid Client Key")
		return
	}

	err := r.ParseForm()
	if err != nil {
		response.BadRequest(w, "ParseForm Invalid form request")
		return
	}

	token := r.PostFormValue("token")
	if token != uuid.Token || token == "" {
		response.BadRequest(w, "Invalid token")
		return
	}

	userid := r.PostFormValue("userid")
	if userid == "" {
		response.BadRequest(w, "UserID is required")
		return
	}

	fundType := r.PostFormValue("fund_type")
	if fundType == "" {
		fundType = "default"
	}

	uid, err := strconv.Atoi(userid)
	if err != nil {
		response.BadRequest(w, "Invalid UserID")
		return
	}

	wallet := NewWallet(h.db)

	if !wallet.IsUserIDUnique(uuid.ID, uid) {
		response.BadRequest(w, "UserID has wallet already")
		return
	}

	wallet.ClientID = uuid.ID
	wallet.UserID = uid
	wallet.Tag = r.PostFormValue("tag")
	wallet.FundType = fundType

	id, err := wallet.Create()
	if err != nil {
		response.InternalServerError(w, err)
		return
	}

	wallet.ID = id
	response.JSON(w, wallet, http.StatusOK)
}

// PostCreditWalletClientUserIDHandle handle the Credit of Client User Wallet
func (h *AppHandler) PostCreditWalletClientUserIDHandle(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	clientUUID := v["clientuuid"]
	uid := v["id"]

	u := NewClient(h.db)
	uuid, Ok := u.ValidUUID(clientUUID)
	if !Ok {
		response.BadRequest(w, "Invalid Client Key")
		return
	}

	if err := r.ParseForm(); err != nil {
		response.BadRequest(w, "Invalid Client Key")
		return
	}
	token := r.PostFormValue("token")
	amount := r.PostFormValue("amount")
	amt, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		response.BadRequest(w, "Invalid Amount Provided")
		return
	}
	CrAmount := float32(amt)

	if token != uuid.Token || token == "" {
		response.BadRequest(w, "Invalid Token")
		return
	}

	walletID, _ := strconv.Atoi(uid)

	wallet := NewWallet(h.db)
	if !wallet.IsClientWalletIDExist(uuid.ID, walletID) {
		response.BadRequest(w, "Wallet not existing")
		return
	}

	t := NewTransaction(h.db)
	t.ClientID = uuid.ID
	t.WalletID = walletID
	t.TransactionType = "cr"
	t.CrAmount = CrAmount
	t.MethodType = r.PostFormValue("method_type")
	t.Particulars = r.PostFormValue("particulars")
	recID, err := t.Credit()
	if err != nil {
		response.BadRequest(w, "Internal server error")
		return
	}
	t.OldBalance, t.NewBalance, err = wallet.CreditWalletBalanceClientWalletID(t.ClientID, t.WalletID, t.CrAmount)
	if err != nil {
		// TODO: revert transactions
		response.BadRequest(w, "Internal server error")
		return
	}

	t.ID = recID
	response.JSON(w, t, http.StatusOK)

	// TODO:
	// 	1) validate the provided ID and Token
	// 	2) check if clientID+UserID has wallet
	//  3) store the API request to the transactions table
	//  4) credit the clientID+UserID wallet
}

// PostDebitWalletClientUserIDHandle handle the Debit of Client User Wallet
func (h *AppHandler) PostDebitWalletClientUserIDHandle(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	clientUUID := v["clientuuid"]
	uid := v["id"]

	u := NewClient(h.db)
	uuid, Ok := u.ValidUUID(clientUUID)
	if !Ok {
		response.BadRequest(w, "Invalid Client Key")
		return
	}

	if err := r.ParseForm(); err != nil {
		response.BadRequest(w, "Invalid Client Key")
		return
	}
	token := r.PostFormValue("token")
	amount := r.PostFormValue("amount")
	amt, err := strconv.ParseFloat(amount, 32)
	if err != nil {
		response.BadRequest(w, "Invalid Amount Provided")
		return
	}
	DrAmount := float32(amt)

	if token != uuid.Token || token == "" {
		response.BadRequest(w, "Invalid Token")
		return
	}

	walletID, _ := strconv.Atoi(uid)

	wallet := NewWallet(h.db)
	if !wallet.IsClientWalletIDExist(uuid.ID, walletID) {
		response.BadRequest(w, "Wallet not existing")
		return
	}

	t := NewTransaction(h.db)
	t.ClientID = uuid.ID
	t.WalletID = walletID
	t.TransactionType = "dr"
	t.DrAmount = DrAmount
	t.MethodType = r.PostFormValue("method_type")
	t.Particulars = r.PostFormValue("particulars")
	recID, err := t.Debit()
	if err != nil {
		response.BadRequest(w, "Internal server error")
		return
	}
	t.OldBalance, t.NewBalance, err = wallet.DebitWalletBalanceClientWalletID(t.ClientID, t.WalletID, t.DrAmount)
	if err != nil {
		// TODO: revert transactions
		response.BadRequest(w, "Internal server error")
		return
	}

	t.ID = recID
	response.JSON(w, t, http.StatusOK)
}

// GetTransactionsClientUserIDHandle returns all transactions of Client User Wallet
func (h *AppHandler) GetTransactionsClientUserIDHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", "GetTransactionsHandler called")
}
