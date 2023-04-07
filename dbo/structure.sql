BEGIN;

-- Create Product table
CREATE TABLE Product (
  id SERIAL PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description TEXT,
  price INTEGER NOT NULL,
  quantity INTEGER NOT NULL
);

-- Create Customer table
CREATE TABLE Customer (
  id SERIAL PRIMARY KEY,
  first_name VARCHAR(50) NOT NULL,
  last_name VARCHAR(50) NOT NULL
);

-- Create Bill table
CREATE TABLE Bill (
  id SERIAL PRIMARY KEY,
  number uuid NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  customer_id INTEGER NOT NULL REFERENCES Customer(id)
);

-- Create ProductBill pivot table
CREATE TABLE ProductBill (
  product_id INTEGER NOT NULL REFERENCES Product(id) ON DELETE CASCADE ON UPDATE CASCADE,
  bill_id INTEGER NOT NULL REFERENCES Bill(id) ON DELETE CASCADE ON UPDATE CASCADE,
  quantity INTEGER NOT NULL,
  PRIMARY KEY (product_id, bill_id)
);

END;