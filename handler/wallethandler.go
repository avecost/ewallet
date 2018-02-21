package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/avecost/ewallet/models"
	"github.com/avecost/ewallet/response"
)

// WalletPostHandler handle the creation of wallet for client subscribers/users
func (h *AppHandler) WalletPostHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	// get client ID from given uuid
	clientID, _ := h.db.GetClientIDByUUID(uuid)
	if clientID == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	// init empty wallet
	wallet := models.Wallet{}
	// decode the pass json object
	err := json.NewDecoder(req.Body).Decode(&wallet)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}
	wallet.ClientID = clientID
	// UserID is required
	if wallet.UserID == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "UserID required"}, http.StatusBadRequest)
		return
	}
	// FundType Default: default
	if wallet.FundType == "" {
		wallet.FundType = "default"
	}
	// Create Wallet
	id, err := h.db.CreateWallet(&wallet)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	newWallet, _ := h.db.GetWalletByID(id)
	response.JSON(w, SuccessResponse{Data: &newWallet}, http.StatusOK)
}

// WalletGetHandler return the Wallet of client subscribers/users
func (h *AppHandler) WalletGetHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	guid := vars["guid"]
	if guid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid wallet"}, http.StatusBadRequest)
		return
	}

	// get client ID from given uuid
	clientID, _ := h.db.GetClientIDByUUID(uuid)
	if clientID == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}

	// get wallet details based on clientID and guid
	wallet, err := h.db.GetWalletByIDGUID(clientID, guid)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	response.JSON(w, SuccessResponse{Data: &wallet}, http.StatusOK)
}

// WalletGetAllHandler returns all Wallet for a client
func (h *AppHandler) WalletGetAllHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}

	// get client ID from given uuid
	clientID, _ := h.db.GetClientIDByUUID(uuid)
	if clientID == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}

	wallets, err := h.db.GetAllWallet(clientID)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	response.JSON(w, SuccessResponse{Data: &wallets}, http.StatusOK)
}

// WalletPutHandler update info of Client Subcribers/Users Wallet
func (h *AppHandler) WalletPutHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}
	guid := vars["guid"]
	if guid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid wallet"}, http.StatusBadRequest)
		return
	}

	// get client ID from given uuid
	clientID, _ := h.db.GetClientIDByUUID(uuid)
	if clientID == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Invalid client uuid"}, http.StatusBadRequest)
		return
	}

	// init empty wallet
	wallet := models.Wallet{}
	// decode the pass json object
	err := json.NewDecoder(req.Body).Decode(&wallet)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}
	// if fundtype not provided make it default
	if wallet.FundType == "" {
		wallet.FundType = "default"
	}
	// check if IsActive is provided if not restore db value
	if wallet.IsActive == nil {
		wallet.IsActive = h.db.GetWalletStatusByIDGUID(clientID, guid)
	}

	// update wallet details based on clientID and guid
	_, err = h.db.UpdateWalletByIDGUID(clientID, guid, &wallet)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	updatedWallet, _ := h.db.GetWalletByIDGUID(clientID, guid)
	response.JSON(w, SuccessResponse{Data: &updatedWallet}, http.StatusOK)
}
