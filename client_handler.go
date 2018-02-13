package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetAllClients return the list of clients
func (h *AppHandler) GetAllClients(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get All Client Handler")
}

// PostNewClient create a new Client
func (h *AppHandler) PostNewClient(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		json.NewEncoder(w).Encode(UrlErrResponse{Version: API_VERSION, Code: http.StatusBadRequest, Message: err.Error()})
	}

	c := &Client{
		db:        h.db,
		Name:      r.PostFormValue("name"),
		Token:     r.PostFormValue("token"),
		Address:   r.PostFormValue("address"),
		URL:       r.PostFormValue("url"),
		Reference: r.PostFormValue("reference"),
	}

	if c.Validate() == false {
		json.NewEncoder(w).Encode(ClientErrorResponse{Version: API_VERSION, Errors: c})
		return
	}

	id, err := c.Create()
	if err != nil {
		appErr := NewErrResponse(API_VERSION)
		appErr.AddError(FieldError{Code: "500", Message: err.Error()})
		appErr.Send(w, http.StatusBadRequest)
		return
	}

	c.ID = id
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ClientSuccessResponse{Version: API_VERSION, Data: c})
}

// GetClientByID return a client record
func (h *AppHandler) GetClientByID(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	id, err := strconv.Atoi(v["id"])
	if err != nil {
		json.NewEncoder(w).Encode(UrlErrResponse{Version: API_VERSION, Code: http.StatusBadRequest, Message: err.Error()})
		return
	}

	c := NewClient(h.db)
	id, _ = c.GetClientById(id)
	if id == 0 {
		w.WriteHeader(http.Status)
		json.NewEncoder(w).Encode(UrlErrResponse{Version: API_VERSION, Code: "", Message: ""})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ClientSuccessResponse{Version: API_VERSION, Data: c})
}

// UpdateClientByID update client record
func (h *AppHandler) UpdateClientByID(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	id, err := strconv.Atoi(v["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Update Client %d Handler", id)
}

// DeleteClientByID mark client record as deleted
func (h *AppHandler) DeleteClientByID(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	id, err := strconv.Atoi(v["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Delete Client %d Handler", id)
}
