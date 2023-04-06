package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func errorHandler(h apiHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			switch t := err.(type) {
			case *ApiError:
				writeJSON(w, http.StatusBadRequest, t)
			default:
				log.Println(err)
				writeJSON(w, http.StatusInternalServerError, nil)
			}
		}
	}
}

type Handler struct {
	s Service
	v *validator.Validate
}

func NewHandler(s Service, v *validator.Validate) *Handler {
	return &Handler{s: s, v: v}
}

func (h *Handler) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/product", errorHandler(h.handleGetProducts)).Methods("GET")
	r.HandleFunc("/product/{id}", errorHandler(h.handleGetProductById)).Methods("GET")
	r.HandleFunc("/product", errorHandler(h.handleAddProduct)).Methods("POST")
	r.HandleFunc("/product/{id}", errorHandler(h.handleUpdateProductById)).Methods("PATCH")
	r.HandleFunc("/product/{id}", errorHandler(h.handleDeleteProductById)).Methods("DELETE")

	r.HandleFunc("/customer", errorHandler(h.handleGetCustomers)).Methods("GET")
	r.HandleFunc("/customer/{id}", errorHandler(h.handleGetCustomerById)).Methods("GET")
	r.HandleFunc("/customer", errorHandler(h.handleAddCustomer)).Methods("POST")
	r.HandleFunc("/customer/{id}", errorHandler(h.handleUpdateCustomerById)).Methods("PATCH")
	r.HandleFunc("/customer/{id}", errorHandler(h.handleDeleteCustomerById)).Methods("DELETE")

	r.HandleFunc("/bill", errorHandler(h.handleGetBills)).Methods("GET")
	r.HandleFunc("/bill/{id}", errorHandler(h.handleGetBillById)).Methods("GET")
	r.HandleFunc("/bill", errorHandler(h.handleAddBill)).Methods("POST")
	r.HandleFunc("/bill/{id}", errorHandler(h.handleUpdateBillById)).Methods("PATCH")
	r.HandleFunc("/bill/{id}", errorHandler(h.handleDeleteBillById)).Methods("DELETE")

	r.HandleFunc("/bill/{id}/product", errorHandler(h.handleGetBillProducts)).Methods("GET")
	r.HandleFunc("/bill/{id}/product", errorHandler(h.handleAddProductToBill)).Methods("POST")
	r.HandleFunc("/bill/{bill_id}/product/{product_id}", errorHandler(h.handleDeleteProductFromBill)).Methods("DELETE")
}

func (h *Handler) handleGetProducts(w http.ResponseWriter, r *http.Request) error {
	products, err := h.s.GetProducts(context.TODO())
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, products)
	return nil
}

func (h *Handler) handleGetProductById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid product's id"}
	}

	product, err := h.s.GetProductById(context.TODO(), id)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, product)
	return nil
}

func (h *Handler) handleAddProduct(w http.ResponseWriter, r *http.Request) error {
	var dto ProductDTOAdd
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		return err
	}

	product, err := h.s.AddProduct(context.TODO(), dto)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, product)
	return nil
}

func (h *Handler) handleUpdateProductById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid product's id"}
	}

	var dto ProductDTOUpdate
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		return err
	}
	dto.Id = id

	if err := h.s.UpdateProductById(context.TODO(), dto); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func (h *Handler) handleDeleteProductById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid product's id"}
	}

	if err := h.s.DeleteProductById(context.TODO(), id); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func (h *Handler) handleGetCustomers(w http.ResponseWriter, r *http.Request) error {
	customers, err := h.s.GetCustomers(context.TODO())
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, customers)
	return nil
}

func (h *Handler) handleGetCustomerById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid customer's id"}
	}

	customer, err := h.s.GetCustomerById(context.TODO(), id)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, customer)
	return nil
}

func (h *Handler) handleAddCustomer(w http.ResponseWriter, r *http.Request) error {
	var dto CustomerDTOAdd
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		return err
	}

	customer, err := h.s.AddCustomer(context.TODO(), dto)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, customer)
	return nil
}

func (h *Handler) handleUpdateCustomerById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid customer's id"}
	}

	var dto CustomerDTOUpdate
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		return err
	}
	dto.Id = id

	if err := h.s.UpdateCustomerById(context.TODO(), dto); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func (h *Handler) handleDeleteCustomerById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid customer's id"}
	}

	if err := h.s.DeleteCustomerById(context.TODO(), id); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func (h *Handler) handleGetBills(w http.ResponseWriter, r *http.Request) error {
	bills, err := h.s.GetBills(context.TODO())
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, bills)
	return nil
}

func (h *Handler) handleGetBillById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid bill's id"}
	}

	bill, err := h.s.GetBillById(context.TODO(), id)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, bill)
	return nil
}

func (h *Handler) handleAddBill(w http.ResponseWriter, r *http.Request) error {
	var dto BillDTOAdd
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		return err
	}

	bill, err := h.s.AddBill(context.TODO(), dto)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, bill)
	return nil
}

func (h *Handler) handleUpdateBillById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid bill's id"}
	}

	var dto BillDTOUpdate
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		return err
	}
	dto.Id = id

	if err := h.s.UpdateBillById(context.TODO(), dto); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func (h *Handler) handleDeleteBillById(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid bill's id"}
	}

	if err := h.s.DeleteBillById(context.TODO(), id); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func (h *Handler) handleAddProductToBill(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid bill's id"}
	}

	var dto BillDtoAddProduct
	if err := decodeAndValidate(&dto, r.Body, h.v, nil); err != nil {
		return err
	}
	dto.Id = id

	err = h.s.AddProductToBill(context.TODO(), dto)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, nil)
	return nil
}

func (h *Handler) handleGetBillProducts(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return &ApiError{Err: "Invalid bill's id"}
	}

	products, err := h.s.GetBillProducts(context.TODO(), id)
	if err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, products)
	return nil
}

func (h *Handler) handleDeleteProductFromBill(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	bill_id, err := strconv.Atoi(vars["bill_id"])
	if err != nil {
		return &ApiError{Err: "Invalid bill's id"}
	}
	product_id, err := strconv.Atoi(vars["product_id"])
	if err != nil {
		return &ApiError{Err: "Invalid product's id"}
	}

	if err := h.s.DeleteProductFromBill(context.TODO(), bill_id, product_id); err != nil {
		return err
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func decodeAndValidate[T any](object T, r io.Reader, v *validator.Validate, addFields func(object T) error) error {
	if err := json.NewDecoder(r).Decode(object); err != nil {
		return &ApiError{Err: "Cannot parse json data"}
	}

	if addFields != nil {
		if err := addFields(object); err != nil {
			return err
		}
	}

	if err := v.Struct(object); err != nil {
		return &ApiError{Err: fmt.Sprintf("Invalid passed data err: %v", err.Error())}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
