package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/avecost/ewallet/models"
	"github.com/avecost/ewallet/response"
	"github.com/gorilla/mux"
)

// CreditPostHandler handle the e-Wallet Credit
func (h *AppHandler) CreditPostHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	// e-Wallet address
	guid := vars["guid"]
	if guid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Client e-Wallet address required"}, http.StatusBadRequest)
		return
	}
	// get client ID from given uuid
	clientID, _ := h.db.GetClientIDByUUID(uuid)
	if clientID == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	// validate if ClientID and e-Wallet address exist
	if !h.db.IsWalletActiveByIDGUID(clientID, guid) {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "e-Wallet not active"}, http.StatusBadRequest)
		return
	}

	// init empty transaction
	creditTransact := models.Transaction{}
	// decode the pass json object
	err := json.NewDecoder(req.Body).Decode(&creditTransact)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}
	creditTransact.ClientID = clientID
	creditTransact.Address = guid
	// validate if crAmount is more than zero
	if creditTransact.CrAmount == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "crAmount must be over zero (0)"}, http.StatusBadRequest)
		return
	}
	// validate if methodType is not empty
	if creditTransact.MethodType == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Method type required"}, http.StatusBadRequest)
		return
	}

	// Create Credit transaction
	id, err := h.db.CreateCreditTransaction(&creditTransact)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Internal server error"}, http.StatusBadRequest)
		return
	}

	// credit the Wallet Balance
	oldBalance, newBalance, err := h.db.CreditWalletByIDGUID(clientID, guid, creditTransact.CrAmount)
	if err != nil {
		fmt.Println(err.Error())
	}
	// get all info about the transaction
	crTransact, _ := h.db.GetTransactionByID(id)
	crTransact.OldBalance = *oldBalance
	crTransact.NewBalance = *newBalance

	response.JSON(w, SuccessResponse{Data: &crTransact}, http.StatusOK)
}

// DebitPostHandler handle the e-Wallet Debit
func (h *AppHandler) DebitPostHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	// e-Wallet address
	guid := vars["guid"]
	if guid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Client e-Wallet address required"}, http.StatusBadRequest)
		return
	}
	// get client ID from given uuid
	clientID, _ := h.db.GetClientIDByUUID(uuid)
	if clientID == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	// validate if ClientID and e-Wallet address exist
	if !h.db.IsWalletActiveByIDGUID(clientID, guid) {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "e-Wallet not active"}, http.StatusBadRequest)
		return
	}

	// init empty transaction
	debitTransact := models.Transaction{}
	// decode the pass json object
	err := json.NewDecoder(req.Body).Decode(&debitTransact)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}
	debitTransact.ClientID = clientID
	debitTransact.Address = guid
	// validate if crAmount is more than zero
	if debitTransact.DrAmount == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "drAmount must be over zero (0)"}, http.StatusBadRequest)
		return
	}
	// validate if methodType is not empty
	if debitTransact.MethodType == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Method type required"}, http.StatusBadRequest)
		return
	}
	// check if has balance for the debit amt
	if !h.db.IsBalanceEnoughForDebit(clientID, guid, debitTransact.DrAmount) {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Insufficient e-Wallet balance"}, http.StatusBadRequest)
		return
	}

	// create Debit transaction
	id, err := h.db.CreateDebitTransaction(&debitTransact)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Internal server error"}, http.StatusBadRequest)
		return
	}

	// debit the Wallet Balance
	oldBalance, newBalance, err := h.db.DebitWalletByIDGUID(clientID, guid, debitTransact.DrAmount)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}
	// get all info about the transaction
	drTransact, _ := h.db.GetTransactionByID(id)
	drTransact.OldBalance = *oldBalance
	drTransact.NewBalance = *newBalance

	response.JSON(w, SuccessResponse{Data: &drTransact}, http.StatusOK)
}

// GetAllTransactionHandler return all transaction of e-Wallet Address
func (h *AppHandler) GetAllTransactionHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	// e-Wallet address
	guid := vars["guid"]
	if guid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Client e-Wallet address required"}, http.StatusBadRequest)
		return
	}
	// get client ID from given uuid
	clientID, _ := h.db.GetClientIDByUUID(uuid)
	if clientID == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	// validate if ClientID and e-Wallet address exist
	if !h.db.IsWalletActiveByIDGUID(clientID, guid) {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "e-Wallet not active"}, http.StatusBadRequest)
		return
	}

	ts, err := h.db.GetAllTransactionByIDGUID(clientID, guid)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	response.JSON(w, SuccessResponse{Data: &ts}, http.StatusOK)
}
