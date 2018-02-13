package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/avecost/ewallet/response"
	"github.com/gorilla/mux"
)

// GetAllClients return the list of clients
func (h *AppHandler) GetAllClients(w http.ResponseWriter, r *http.Request) {
	//TODO: return all clients
	fmt.Fprintf(w, "Get All Client Handler")
}

// PostNewClient create a new Client
func (h *AppHandler) PostNewClient(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.BadRequest(w, "Invalid form request")
		return
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

// GetClientByID return a client record
func (h *AppHandler) GetClientByID(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	id, err := strconv.Atoi(v["id"])
	if err != nil {
		response.BadRequest(w, "Incorrect request")
		return
	}

	c := NewClient(h.db)
	id, _ = c.GetClientByID(id)

	response.JSON(w, c, http.StatusOK)
}

// UpdateClientByID update client record
func (h *AppHandler) UpdateClientByID(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	id, err := strconv.Atoi(v["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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
		response.BadRequest(w, c)
		return
	}

	err = c.UpdateClientByID(id)
	if err != nil {
		response.InternalServerError(w, "Error with request")
		return
	}

	c.ID = id
	response.JSON(w, c, http.StatusOK)
}

// DeleteClientByID mark client record as deleted
func (h *AppHandler) DeleteClientByID(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)

	id, err := strconv.Atoi(v["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := NewClient(h.db)
	err = c.DeleteClientByID(id)
	if err != nil {
		response.InternalServerError(w, "Error with request")
		return
	}

	c.ID = id
	response.JSON(w, c, http.StatusOK)
}
