package main

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type Handler struct {
	s Service
	v *validator.Validate
}

func NewHandler(s Service, v *validator.Validate) *Handler {
	return &Handler{s: s, v: v}
}

func (h *Handler) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/product", h.passHandle).Methods("GET")
	r.HandleFunc("/product/{id}", h.passHandle).Methods("GET")
	r.HandleFunc("/product", h.passHandle).Methods("POST")
	r.HandleFunc("/product/{id}", h.passHandle).Methods("PATCH")
	r.HandleFunc("/product/{id}", h.passHandle).Methods("DELETE")

	r.HandleFunc("/customer", h.passHandle).Methods("GET")
	r.HandleFunc("/customer/{id}", h.passHandle).Methods("GET")
	r.HandleFunc("/customer", h.passHandle).Methods("POST")
	r.HandleFunc("/customer/{id}", h.passHandle).Methods("PATCH")
	r.HandleFunc("/customer/{id}", h.passHandle).Methods("DELETE")

	r.HandleFunc("/bill", h.passHandle).Methods("GET")
	r.HandleFunc("/bill/{id}", h.passHandle).Methods("GET")
	r.HandleFunc("/bill", h.passHandle).Methods("POST")
	r.HandleFunc("/bill/{id}", h.passHandle).Methods("PATCH")
	r.HandleFunc("/bill/{id}", h.passHandle).Methods("DELETE")

	r.HandleFunc("/bill/{id}/product", h.passHandle).Methods("GET")
	r.HandleFunc("/bill/{id}/product", h.passHandle).Methods("POST")
	r.HandleFunc("/bill/{id}/product/{id}", h.passHandle).Methods("DELETE")
}

func (h *Handler) passHandle(w http.ResponseWriter, r *http.Request) {
	panic("Not implemented")
}
