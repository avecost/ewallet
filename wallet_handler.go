package main

import (
	"net/http"
	"strconv"

	"github.com/avecost/ewallet/response"
)

// PostNewWallet create a new Wallet
func (h *AppHandler) PostNewWallet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.BadRequest(w, "Invalid form request")
		return
	}

	cID, err := strconv.Atoi(r.PostFormValue("client_id"))
	if err != nil {
		response.BadRequest(w, "Invalid form request")
		return
	}
	uID, err := strconv.Atoi(r.PostFormValue("user_id"))
	if err != nil {
		response.BadRequest(w, "Invalid form request")
		return
	}
	balance, err := strconv.ParseFloat(r.PostFormValue("balance"), 64)
	if err != nil {
		response.BadRequest(w, "Invalid form request")
		return
	}

	w := NewWallet(h.db)
	w.ClientID = cID
	w.UserID = uID
	w.Tag = r.PostFormValue("tag")
	w.Balance = balance
	w.FundType = r.PostFormValue("fund_type")

	

	if c.Validate() == false {
		response.BadRequest(w, c)
		return
	}

	id, err := c.Create()
	if err != nil {
		response.InternalServerError(w, "Error with request")
		return
	}

	c.ID = id
	w.WriteHeader(http.StatusOK)
	response.JSON(w, c, http.StatusOK)
	// json.NewEncoder(w).Encode(ClientSuccessResponse{Version: API_VERSION, Data: c})
}
