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
