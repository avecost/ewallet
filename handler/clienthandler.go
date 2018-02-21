package handler

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/avecost/ewallet/models"
	"github.com/avecost/ewallet/response"
)

// ClientPostHandler handle the new client
func (h *AppHandler) ClientPostHandler(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	name := req.PostFormValue("name")
	if name == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Name required"}, http.StatusBadRequest)
		return
	}
	token := req.PostFormValue("token")
	if token == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Token required"}, http.StatusBadRequest)
		return
	}

	c := &models.Client{
		Name:      name,
		Token:     token,
		Address:   req.PostFormValue("address"),
		URL:       req.PostFormValue("url"),
		Reference: req.PostFormValue("reference"),
	}

	id, err := h.db.CreateClient(c)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	client, _ := h.db.GetClientByID(id)
	response.JSON(w, SuccessResponse{Data: &client}, http.StatusOK)
}

// ClientGetHandler return client record based on uuid
func (h *AppHandler) ClientGetHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	uuid := v["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Client UUID required"}, http.StatusBadRequest)
		return
	}
	client, err := h.db.GetClientByUUID(uuid)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	response.JSON(w, SuccessResponse{Data: &client}, http.StatusOK)
}

// ClientGetAllHandler returns array of client objects
func (h *AppHandler) ClientGetAllHandler(w http.ResponseWriter, req *http.Request) {
	clients, err := h.db.GetAllClient()
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	response.JSON(w, SuccessResponse{Data: &clients}, http.StatusOK)
}

// ClientPutHandler update client object
func (h *AppHandler) ClientPutHandler(w http.ResponseWriter, req *http.Request) {
	v := mux.Vars(req)
	uuid := v["uuid"]
	if uuid == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Client UUID required"}, http.StatusBadRequest)
		return
	}

	err := req.ParseForm()
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}

	c := &models.Client{}

	name := req.PostFormValue("name")
	if name == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Name required"}, http.StatusBadRequest)
		return
	}
	c.Name = name

	token := req.PostFormValue("token")
	if token == "" {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Token required"}, http.StatusBadRequest)
		return
	}
	c.Token = token

	address := req.PostFormValue("address")
	if address != "" {
		c.Address = address
	}

	url := req.PostFormValue("url")
	if url != "" {
		c.URL = url
	}

	reference := req.PostFormValue("reference")
	if reference != "" {
		c.Reference = reference
	}

	res, err := h.db.UpdateClientByUUID(uuid, c)
	if err != nil {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: err.Error()}, http.StatusBadRequest)
		return
	}
	if res == 0 {
		response.JSON(w, ErrResponse{Err: "Application Error", Message: "Record not found"}, http.StatusBadRequest)
		return
	}

	client, _ := h.db.GetClientByUUID(uuid)
	response.JSON(w, SuccessResponse{Data: &client}, http.StatusOK)
}
