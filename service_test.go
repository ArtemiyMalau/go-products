package main

import (
	"context"
	"log"
	"sync"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type Environment struct {
	s Service
}

var onceTest sync.Once
var environmentInstance *Environment

func GetEnvironment() *Environment {
	onceTest.Do(func() {
		config := GetConfig()

		db, err := GetDBConnection(config)
		if err != nil {
			log.Fatal(err)
		}

		service := NewService(db)
		environmentInstance = &Environment{s: *service}
	})
	return environmentInstance
}

func TestProduct(t *testing.T) {
	e := GetEnvironment()
	dtoAdd := ProductDTOAdd{
		Name:        "Test Product",
		Description: "Description of test product",
		Price:       100000,
		Quantity:    100000,
	}
	var product *Product
	var err error

	// Create product
	product, err = e.s.AddProduct(context.TODO(), dtoAdd)
	if err != nil {
		t.Errorf("Error when adding product: %+v", err)
	}
	if product.Name != dtoAdd.Name ||
		product.Description != dtoAdd.Description ||
		product.Price != dtoAdd.Price ||
		product.Quantity != dtoAdd.Quantity {
		t.Errorf("Added product have different values: dto:%+v; product:%+v", dtoAdd, product)
	}

	// Get product by ID
	var productCopy *Product
	productCopy, err = e.s.GetProductById(context.TODO(), product.Id)
	if err != nil {
		t.Errorf("Error when fetching product: %+v", err)
	}
	if !cmp.Equal(product, productCopy) {
		t.Errorf("Error when compare inserted and fetched product: %+v", err)
	}

	// Get products (checking whether created product in products slice)
	var products []Product
	products, err = e.s.GetProducts(context.TODO())
	if err != nil {
		t.Errorf("Error when fetching products: %+v", err)
	}

	isHaveInProducts := false
	for _, p := range products {
		if product.Id == p.Id {
			isHaveInProducts = true
		}
	}
	if !isHaveInProducts {
		t.Errorf("Created product are not in products")
	}

	// UpdateProduct
	dtoUpdate := ProductDTOUpdate{
		Id:          product.Id,
		Name:        "Updated Test Product",
		Description: "Updated Description of test product",
		Price:       50000,
		Quantity:    50000,
	}
	err = e.s.UpdateProductById(context.TODO(), dtoUpdate)
	if err != nil {
		t.Errorf("Error when updating product: %+v", err)
	}

	// Delete product
	err = e.s.DeleteProductById(context.TODO(), product.Id)
	if err != nil {
		t.Errorf("Error when deleting product: %+v", err)
	}

	var undesiderProduct *Product
	undesiderProduct, err = e.s.GetProductById(context.TODO(), product.Id)
	if undesiderProduct != nil {
		t.Errorf("Get product by id return not nil product after deleting: %+v", undesiderProduct)
	}
	if err != nil {
		if _, ok := err.(*ApiError); !ok {
			t.Errorf("Error when deleting product: %+v", err)
		}
	} else {
		t.Errorf("Get product must return not found ApiError")
	}
}
