package main

import (
	"context"
	"fmt"
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
			t.Errorf("Error when fetching product: %+v", err)
		}
	} else {
		t.Errorf("Get product must return not found ApiError")
	}
}

func TestCustomer(t *testing.T) {
	e := GetEnvironment()
	dtoAdd := CustomerDTOAdd{
		FirstName: "Test name",
		LastName:  "Test last name",
	}
	var customer *Customer
	var err error

	// Create customer
	customer, err = e.s.AddCustomer(context.TODO(), dtoAdd)
	if err != nil {
		t.Errorf("Error when adding customer: %+v", err)
	}
	if customer.FirstName != dtoAdd.FirstName ||
		customer.LastName != dtoAdd.LastName {
		t.Errorf("Added customer have different values: dto:%+v; customer:%+v", dtoAdd, customer)
	}

	// Get customer by ID
	var customerCopy *Customer
	customerCopy, err = e.s.GetCustomerById(context.TODO(), customer.Id)
	if err != nil {
		t.Errorf("Error when fetching customer: %+v", err)
	}
	if !cmp.Equal(customer, customerCopy) {
		t.Errorf("Error when compare inserted and fetched customers: %+v", err)
	}

	// Get customers (checking whether created customer in customers slice)
	var customers []Customer
	customers, err = e.s.GetCustomers(context.TODO())
	if err != nil {
		t.Errorf("Error when fetching customers: %+v", err)
	}

	isHaveInCustomers := false
	for _, c := range customers {
		if customer.Id == c.Id {
			isHaveInCustomers = true
		}
	}
	if !isHaveInCustomers {
		t.Errorf("Created customer are not in customers")
	}

	// UpdateCustomer
	dtoUpdate := CustomerDTOUpdate{
		Id:        customer.Id,
		FirstName: "Updated Test Name",
		LastName:  "Updated Test Last Name",
	}
	err = e.s.UpdateCustomerById(context.TODO(), dtoUpdate)
	if err != nil {
		t.Errorf("Error when updating customer: %+v", err)
	}

	// Delete customer
	err = e.s.DeleteCustomerById(context.TODO(), customer.Id)
	if err != nil {
		t.Errorf("Error when deleting customer: %+v", err)
	}

	var undesiderCustomer *Customer
	undesiderCustomer, err = e.s.GetCustomerById(context.TODO(), customer.Id)
	if undesiderCustomer != nil {
		t.Errorf("Get customer by id return not nil customer after deleting: %+v", undesiderCustomer)
	}
	if err != nil {
		if _, ok := err.(*ApiError); !ok {
			t.Errorf("Error when fetching customer: %+v", err)
		}
	} else {
		t.Errorf("Get product must return not found ApiError")
	}
}

func TestBill(t *testing.T) {
	e := GetEnvironment()

	// setup
	var (
		productOne, productTwo   *Product
		customerOne, customerTwo *Customer
	)
	err := func() error {
		var err error
		productOne, err = e.s.AddProduct(context.TODO(), ProductDTOAdd{
			Name:        "Product One",
			Description: "Description",
			Price:       1000,
			Quantity:    10,
		})
		if err != nil {
			return err
		}
		productTwo, err = e.s.AddProduct(context.TODO(), ProductDTOAdd{
			Name:        "Product One",
			Description: "Description",
			Price:       1000,
			Quantity:    10,
		})
		if err != nil {
			return err
		}
		customerOne, err = e.s.AddCustomer(context.TODO(), CustomerDTOAdd{
			FirstName: "Customer One",
			LastName:  "Customer One",
		})
		if err != nil {
			return err
		}
		customerTwo, err = e.s.AddCustomer(context.TODO(), CustomerDTOAdd{
			FirstName: "Customer Two",
			LastName:  "Customer Two",
		})
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		panic(fmt.Sprintf("Cannot setup TestBill: %+v", err))
	}
	// end setup

	dtoAdd := BillDTOAdd{
		Customer: customerOne.Id,
		Products: []BillProduct{
			{
				Product:  productOne.Id,
				Quantity: 10,
			},
			{
				Product:  productTwo.Id,
				Quantity: 5,
			},
		},
	}

	// Create bill
	bill, err := e.s.AddBill(context.TODO(), dtoAdd)
	if err != nil {
		t.Errorf("Error when adding bill: %+v", err)
	}
	if bill.Customer != dtoAdd.Customer {
		t.Errorf("Added customer have different values: dto:%+v; customer:%+v", dtoAdd, bill)
	}

	// Get bills (checking whether created bill in bills slice)
	var bills []Bill
	bills, err = e.s.GetBills(context.TODO())
	if err != nil {
		t.Errorf("Error when fetching bills: %+v", err)
	}

	isHaveInBills := false
	for _, b := range bills {
		if bill.Id == b.Id {
			isHaveInBills = true
		}
	}
	if !isHaveInBills {
		t.Errorf("Created bill are not in bills")
	}

	// Get bill verbose by ID
	var billVerbose *BillVerbose
	billVerbose, err = e.s.GetBillById(context.TODO(), bill.Id)
	if err != nil {
		t.Errorf("Error when fetching bill: %+v", err)
	}

	// UpdateBill
	dtoUpdate := BillDTOUpdate{
		Id:       bill.Id,
		Customer: customerTwo.Id,
		Products: []BillProduct{
			{
				Product:  productOne.Id,
				Quantity: 10,
			},
		},
	}
	err = e.s.UpdateBillById(context.TODO(), dtoUpdate)
	if err != nil {
		t.Errorf("Error when updating bill: %+v", err)
	}

	// Get bill verbose by ID
	var billVerboseCopy *BillVerbose
	billVerboseCopy, err = e.s.GetBillById(context.TODO(), bill.Id)
	if err != nil {
		t.Errorf("Error when fetching bill: %+v", err)
	}

	if cmp.Equal(billVerbose, billVerboseCopy) {
		t.Errorf("Created bill and updated bill are same:\n created: %+v\n updated: %+v", billVerbose, billVerboseCopy)
	}

	// Delete bill
	err = e.s.DeleteBillById(context.TODO(), bill.Id)
	if err != nil {
		t.Errorf("Error when deleting bill: %+v", err)
	}

	var undesiredBillVerbose *BillVerbose
	undesiredBillVerbose, err = e.s.GetBillById(context.TODO(), bill.Id)
	if undesiredBillVerbose != nil {
		t.Errorf("Get bill by id return not nil bill after deleting: %+v", undesiredBillVerbose)
	}
	if err != nil {
		if _, ok := err.(*ApiError); !ok {
			t.Errorf("Error when fetching bill: %+v", err)
		}
	} else {
		t.Errorf("Get product must return not found ApiError")
	}

	// teardown
	err = func() error {
		var err error
		err = e.s.DeleteProductById(context.TODO(), productOne.Id)
		if err != nil {
			return err
		}
		err = e.s.DeleteProductById(context.TODO(), productTwo.Id)
		if err != nil {
			return err
		}
		err = e.s.DeleteCustomerById(context.TODO(), customerOne.Id)
		if err != nil {
			return err
		}
		err = e.s.DeleteCustomerById(context.TODO(), customerTwo.Id)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		panic(fmt.Sprintf("Cannot teardown TestBill: %+v", err))
	}
	// end teardown
}

func TestBillProduct(t *testing.T) {
	e := GetEnvironment()

	// setup
	var products []Product
	var bill *Bill
	err := func() error {
		var err error
		products, err = e.s.GetProducts(context.TODO())
		if err != nil {
			return err
		}
		if len(products) < 2 {
			return fmt.Errorf("Count of products have to be at least 2")
		}

		var customers []Customer
		customers, err = e.s.GetCustomers(context.TODO())
		if len(customers) < 1 {
			return fmt.Errorf("Count of customers have to be at least 1")
		}

		bill, err = e.s.AddBill(context.TODO(), BillDTOAdd{
			Customer: customers[0].Id,
			Products: []BillProduct{
				{
					Product:  products[0].Id,
					Quantity: 10,
				},
			},
		})
		if err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		panic(fmt.Sprintf("Cannot setup TestBillProduct: %+v", err))
	}
	// end setup

	billProducts, err := e.s.GetBillProducts(context.TODO(), bill.Id)
	if err != nil {
		t.Errorf("Error when fetching bill's products: %+v", err)
	}
	if len(billProducts) != 1 {
		t.Errorf("Invalid bill products count: have to be 1, got %d", len(billProducts))
	}

	err = e.s.AddProductToBill(context.TODO(), BillDtoAddProduct{
		Id: bill.Id,
		BillProduct: BillProduct{
			Product:  products[1].Id,
			Quantity: 5,
		},
	})
	if err != nil {
		t.Errorf("Error when add product to bill: %+v", err)
	}

	err = e.s.DeleteProductFromBill(context.TODO(), bill.Id, products[1].Id)
	if err != nil {
		t.Errorf("Error when deleting product to bill: %+v", err)
	}

	// teardown
	if err := e.s.DeleteBillById(context.TODO(), bill.Id); err != nil {
		panic(fmt.Sprintf("Cannot teardown TestBillProduct: %+v", err))
	}
	// end teardown
}
