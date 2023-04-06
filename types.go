package main

import (
	"time"

	"github.com/google/uuid"
)

// Product-related types
type Product struct {
	Id          int    `json:"id" db="id"`
	Name        string `json:"name" db="name"`
	Description string `json:"description" db="description"`
	Price       int    `json:"price" db="price"`
	Quantity    int    `json:"quantity" db="quantity"`
}

type ProductDTOAdd struct {
	Name        string `json:"name" validate:"required" db="name"`
	Description string `json:"description" validate:"required" db="description"`
	Price       int    `json:"price" validate:"required,gt=0" db="price"`
	Quantity    int    `json:"quantity" validate:"required,gt=0" db="quantity"`
}

type ProductDTOUpdate struct {
	Id          int    `db:"id"`
	Name        string `json:"name" validate:"required" db="name"`
	Description string `json:"description" validate:"required" db="description"`
	Price       int    `json:"price" validate:"required,gt=0" db="price"`
	Quantity    int    `json:"quantity" validate:"required,gt=0" db="quantity"`
}

// Customer-related types
type Customer struct {
	Id        int    `json:"id" db="id"`
	FirstName string `json:"first_name" db="first_name"`
	LastName  string `json:"last_name" db="last_name"`
}

type CustomerDTOAdd struct {
	FirstName string `json:"first_name" validate:"required" db="first_name"`
	LastName  string `json:"last_name" validate:"required" db="last_name"`
}

type CustomerDTOUpdate struct {
	Id        int    `db:"id"`
	FirstName string `json:"first_name" validate:"required" db="first_name"`
	LastName  string `json:"last_name" validate:"required" db="last_name"`
}

// Bill-related types
type Bill struct {
	Id        int       `json:"id" db:"id"`
	Number    uuid.UUID `json:"number" db:"number"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Customer  int       `json:"customer" db:"customer_id"`
}

type BillVerbose struct {
	Id        int           `json:"id"`
	Number    uuid.UUID     `json:"number"`
	CreatedAt time.Time     `json:"created_at"`
	Customer  Customer      `json:"customer"`
	Products  []BillProduct `json: "products"`
}

type BillProduct struct {
	Product  int `json:"product" validate:"required" db:"product"`
	Quantity int `json:"quantity" validate:"required,gt=0" db:"quantity"`
}

type BillDTOAdd struct {
	Customer int           `json:"customer" validate:"required"`
	Products []BillProduct `json:"products" validate:"required,dive"`
}

type BillDTOUpdate struct {
	Id       int           `db:"id"`
	Customer int           `json:"customer" validate:"required"`
	Products []BillProduct `json:"products" validate:"required,dive"`
}

type BillDtoAddProduct struct {
	BillProduct
	Id int `db:"id"`
}
