package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

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
	r.HandleFunc("/product", h.handleGetProducts).Methods("GET")
	r.HandleFunc("/product/{id}", h.handleGetProductById).Methods("GET")
	r.HandleFunc("/product", h.handleAddProduct).Methods("POST")
	r.HandleFunc("/product/{id}", h.handleUpdateProductById).Methods("PATCH")
	r.HandleFunc("/product/{id}", h.handleDeleteProductById).Methods("DELETE")

	r.HandleFunc("/customer", h.handleGetCustomers).Methods("GET")
	r.HandleFunc("/customer/{id}", h.handleGetCustomerById).Methods("GET")
	r.HandleFunc("/customer", h.handleAddCustomer).Methods("POST")
	r.HandleFunc("/customer/{id}", h.handleUpdateCustomerById).Methods("PATCH")
	r.HandleFunc("/customer/{id}", h.handleDeleteCustomerById).Methods("DELETE")

	r.HandleFunc("/bill", h.handleGetBills).Methods("GET")
	r.HandleFunc("/bill/{id}", h.handleGetBillById).Methods("GET")
	r.HandleFunc("/bill", h.handleAddBill).Methods("POST")
	r.HandleFunc("/bill/{id}", h.handleUpdateBillById).Methods("PATCH")
	r.HandleFunc("/bill/{id}", h.handleDeleteBillById).Methods("DELETE")

	r.HandleFunc("/bill/{id}/product", h.handleGetBillProducts).Methods("GET")
	r.HandleFunc("/bill/{id}/product", h.handleAddProductToBill).Methods("POST")
	r.HandleFunc("/bill/{bill_id}/product/{product_id}", h.handleDeleteProductFromBill).Methods("DELETE")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.s.GetProducts(context.TODO())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, products)
}

func (h *Handler) handleGetProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	product, err := h.s.GetProductById(context.TODO(), id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, product)
}

func (h *Handler) handleAddProduct(w http.ResponseWriter, r *http.Request) {
	var dto ProductDTOAdd
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.s.AddProduct(context.TODO(), dto); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUpdateProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var dto ProductDTOUpdate
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	dto.Id = id

	if err := h.s.UpdateProductById(context.TODO(), dto); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleDeleteProductById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.s.DeleteProductById(context.TODO(), id); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleGetCustomers(w http.ResponseWriter, r *http.Request) {
	customers, err := h.s.GetCustomers(context.TODO())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, customers)
}

func (h *Handler) handleGetCustomerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	customer, err := h.s.GetCustomerById(context.TODO(), id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, customer)
}

func (h *Handler) handleAddCustomer(w http.ResponseWriter, r *http.Request) {
	var dto CustomerDTOAdd
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.s.AddCustomer(context.TODO(), dto); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUpdateCustomerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var dto CustomerDTOUpdate
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	dto.Id = id

	if err := h.s.UpdateCustomerById(context.TODO(), dto); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleDeleteCustomerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.s.DeleteCustomerById(context.TODO(), id); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleGetBills(w http.ResponseWriter, r *http.Request) {
	bills, err := h.s.GetBills(context.TODO())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, bills)
}

func (h *Handler) handleGetBillById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	bill, err := h.s.GetBillById(context.TODO(), id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, bill)
}

func (h *Handler) handleAddBill(w http.ResponseWriter, r *http.Request) {
	var dto BillDTOAdd
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.s.AddBill(context.TODO(), dto); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleUpdateBillById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var dto BillDTOUpdate
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	dto.Id = id

	if err := h.s.UpdateBillById(context.TODO(), dto); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleDeleteBillById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.s.DeleteBillById(context.TODO(), id); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleAddProductToBill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	var dto BillDtoAddProduct
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	dto.Id = id

	err = h.s.AddProductToBill(context.TODO(), dto)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleGetBillProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	products, err := h.s.GetBillProducts(context.TODO(), id)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, products)
}

func (h *Handler) handleDeleteProductFromBill(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bill_id, err := strconv.Atoi(vars["bill_id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	product_id, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.s.DeleteProductFromBill(context.TODO(), bill_id, product_id); err != nil {
		writeJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, nil)
}

func decodeAndValidate[T any](object T, r io.Reader, v *validator.Validate, addFields func(object T) error) error {
	if err := json.NewDecoder(r).Decode(object); err != nil {
		return err
	}

	if addFields != nil {
		if err := addFields(object); err != nil {
			return err
		}
	}

	if err := v.Struct(object); err != nil {
		return err
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
