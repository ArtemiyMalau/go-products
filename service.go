package main

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{
		db: db,
	}
}

// Product-related methods
func (s *Service) GetProducts(ctx context.Context) (products []Product, err error) {
	products = []Product{}
	if err = s.db.SelectContext(ctx, &products, `
	SELECT product.id, product.name, product.description, product.price, product.quantity FROM product
	`); err != nil {
		return
	}
	return
}

func (s *Service) GetProductById(ctx context.Context, id int) (product Product, err error) {
	if err = s.db.GetContext(ctx, &product, `
	SELECT product.id, product.name, product.description, product.price, product.quantity FROM product
	WHERE id = $1
	`, id); err != nil {
		return
	}
	return
}

func (s *Service) AddProduct(ctx context.Context, dto ProductDTOAdd) error {
	if _, err := s.db.NamedExecContext(ctx, `
	INSERT INTO product (name, description, price, quantity) VALUES (:name, :description, :price, :quantity)
	`, &dto); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateProductById(ctx context.Context, dto ProductDTOUpdate) error {
	if _, err := s.db.NamedExecContext(ctx, `
	UPDATE product
	SET name = :name, description = :description, price = :price, quantity = :quantity
	WHERE id = :id
	`, &dto); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteProductById(ctx context.Context, id int) error {
	if _, err := s.db.ExecContext(ctx, `
	DELETE FROM product WHERE id = $1
	`, id); err != nil {
		return err
	}
	return nil
}

// Customer-related methods
func (s *Service) GetCustomers(ctx context.Context) (customers []Customer, err error) {
	customers = []Customer{}
	if err = s.db.SelectContext(ctx, &customers, `
	SELECT customer.id, customer.first_name, customer.last_name FROM customer
	`); err != nil {
		return
	}
	return
}

func (s *Service) GetCustomerById(ctx context.Context, id int) (customer Customer, err error) {
	if err = s.db.GetContext(ctx, &customer, `
	SELECT customer.id, customer.first_name, customer.last_name FROM customer
	WHERE id = $1
	`, id); err != nil {
		return
	}
	return
}

func (s *Service) AddCustomer(ctx context.Context, dto CustomerDTOAdd) error {
	if _, err := s.db.NamedExecContext(ctx, `
	INSERT INTO customer (first_name, last_name) VALUES (:first_name, :last_name)
	`, &dto); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateCustomerById(ctx context.Context, dto CustomerDTOUpdate) error {
	if _, err := s.db.NamedExecContext(ctx, `
	UPDATE product
	SET first_name = :first_name, last_name = :last_name
	WHERE id = :id
	`, &dto); err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteCustomerById(ctx context.Context, id int) error {
	if _, err := s.db.ExecContext(ctx, `
	DELETE FROM customer WHERE id = $1
	`, id); err != nil {
		return err
	}
	return nil
}

// Bill-related methods
func (s *Service) GetBills(ctx context.Context) (bills []Bill, err error) {
	bills = []Bill{}
	if err = s.db.SelectContext(ctx, &bills, `
	SELECT bill.id, bill.number, bill.created_at, bill.customer_id FROM bill
	`); err != nil {
		return
	}
	return
}

func (s *Service) GetBillById(ctx context.Context, id int) (bill BillVerbose, err error) {
	tx := s.db.MustBeginTx(ctx, nil)
	err = func() error {
		err := tx.QueryRowContext(ctx, `
		SELECT bill.id AS bill_id, bill.number, bill.created_at, customer.id AS customer_id, customer.first_name, customer.last_name
		FROM bill
		JOIN customer ON bill.customer_id = customer.id
		WHERE bill.id = $1
		`, id).Scan(
			&bill.Id,
			&bill.Number,
			&bill.CreatedAt,
			&bill.Customer.Id,
			&bill.Customer.FirstName,
			&bill.Customer.LastName,
		)
		if err != nil {
			return err
		}

		if err := tx.SelectContext(ctx, &bill.Products, `
		SELECT productbill.product_id AS product, productbill.quantity FROM productbill WHERE productbill.bill_id = $1
		`, id); err != nil {
			return err
		}

		return nil
	}()
	if err != nil {
		return
	}
	return
}

func (s *Service) validateUpsertBillFields(ctx context.Context, tx *sqlx.Tx, customer int, products []BillProduct) error {
	// Checking customer on existence
	var hasCustomer bool
	tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT id FROM customer WHERE id = $1)", customer).Scan(&hasCustomer)
	if !hasCustomer {
		return fmt.Errorf("customer with passed id:%v not exists", customer)
	}

	// Checking that all passed products have to exists
	buf := bytes.NewBufferString("SELECT COUNT(*) FROM product WHERE id IN (")
	for i, billProduct := range products {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.Itoa(billProduct.Product))
	}
	buf.WriteString(")")

	var realProductsCount int
	tx.QueryRowContext(ctx, buf.String()).Scan(&realProductsCount)
	if len(products) != realProductsCount {
		return fmt.Errorf("not all passed products exists")
	}

	return nil
}

func (s *Service) AddBill(ctx context.Context, dto BillDTOAdd) error {
	tx := s.db.MustBeginTx(ctx, nil)
	if err := func() error {
		if err := s.validateUpsertBillFields(ctx, tx, dto.Customer, dto.Products); err != nil {
			return err
		}

		var billId int
		err := tx.QueryRowContext(ctx, `
		INSERT INTO bill (number, customer_id) VALUES ($1, $2) RETURNING id
		`, uuid.New().String(), dto.Customer).Scan(&billId)
		if err != nil {
			return err
		}

		for _, billProduct := range dto.Products {
			if _, err := tx.ExecContext(ctx, `
			INSERT INTO productbill (product_id, bill_id, quantity) VALUES ($1, $2, $3)
			`, billProduct.Product, billId, billProduct.Quantity); err != nil {
				return err
			}
		}

		return nil
	}(); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *Service) UpdateBillById(ctx context.Context, dto BillDTOUpdate) error {
	tx := s.db.MustBeginTx(ctx, nil)
	if err := func() error {
		var hasBill bool
		tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT id FROM customer WHERE id = $1)", dto.Customer).Scan(&hasBill)
		if !hasBill {
			return fmt.Errorf("Bill with passed id:%v not exists", dto.Id)
		}

		if err := s.validateUpsertBillFields(ctx, tx, dto.Customer, dto.Products); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, `
		UPDATE bill SET customer_id = $1 WHERE id = $2
		`, dto.Customer, dto.Id); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, "DELETE FROM productbill WHERE bill_id = $1", dto.Id); err != nil {
			return err
		}

		for _, billProduct := range dto.Products {
			if _, err := tx.ExecContext(ctx, `
			INSERT INTO productbill (product_id, bill_id, quantity) VALUES ($1, $2, $3)
			`, billProduct.Product, dto.Id, billProduct.Quantity); err != nil {
				return err
			}
		}

		return nil
	}(); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *Service) DeleteBillById(ctx context.Context, id int) error {
	if _, err := s.db.ExecContext(ctx, `
	DELETE FROM bill WHERE id = $1
	`, id); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetBillProducts(ctx context.Context, id int) (products []Product, err error) {
	products = []Product{}
	if err = s.db.SelectContext(ctx, &products, `
	SELECT product.id, product.name, product.description, product.price, product.quantity
	FROM product
	JOIN productbill ON productbill.product_id = product.id
	WHERE productbill.bill_id = $1
	`, id); err != nil {
		return
	}
	return
}

func (s *Service) DeleteProductFromBill(ctx context.Context, bill_id, product_id int) error {
	if _, err := s.db.ExecContext(ctx, `
	DELETE FROM productbill WHERE bill_id = $1 and product_id = $2
	`, bill_id, product_id); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddProductToBill(ctx context.Context, dto BillDtoAddProduct) error {
	if _, err := s.db.ExecContext(ctx, `
	INSERT INTO productbill (product_id, bill_id, quantity) VALUES ($1, $2, $3)
	`, dto.BillProduct.Product, dto.Id, dto.BillProduct.Quantity); err != nil {
		return err
	}
	return nil
}
