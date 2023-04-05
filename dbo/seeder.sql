-- Insert default data for Product table
INSERT INTO Product (name, description, price, quantity) VALUES
('iPhone X', 'Apple iPhone X with OLED screen and Face ID', 999, 50),
('Samsung Galaxy S9', 'Samsung Galaxy S9 with Infinity Display and Bixby', 799, 75),
('MacBook Pro', 'Apple MacBook Pro with Retina Display and Touch Bar', 1999, 25),
('Dell XPS 13', 'Dell XPS 13 with InfinityEdge Display and Windows 10', 1299, 30),
('Nintendo Switch', 'Nintendo Switch with Joy-Con controllers and dock', 299, 100),
('PlayStation 4 Pro', 'Sony PlayStation 4 Pro with 4K HDR gaming and streaming', 399, 80),
('Xbox One X', 'Microsoft Xbox One X with 4K gaming and UHD Blu-ray player', 499, 70),
('Amazon Echo Dot', 'Amazon Echo Dot with Alexa voice assistant and smart home control', 49, 150),
('Google Home Mini', 'Google Home Mini with Google Assistant and Chromecast support', 39, 200),
('Fitbit Charge 2', 'Fitbit Charge 2 with heart rate monitor and fitness tracking', 149, 125);

-- Insert default data for Customer table
INSERT INTO Customer (first_name, last_name) VALUES
('John', 'Doe'),
('Jane', 'Doe'),
('Bob', 'Smith'),
('Alice', 'Johnson'),
('Tom', 'Jones');

-- Insert default data for Bill table
INSERT INTO Bill (number, customer_id) VALUES
('d511023c-d3a9-11ed-afa1-0242ac120002', 1),
('d5110692-d3a9-11ed-afa1-0242ac120002', 2),
('d511087c-d3a9-11ed-afa1-0242ac120002', 3),
('d5110a34-d3a9-11ed-afa1-0242ac120002', 4),
('d511102e-d3a9-11ed-afa1-0242ac120002', 5),
('d5111240-d3a9-11ed-afa1-0242ac120002', 1),
('d51114e8-d3a9-11ed-afa1-0242ac120002', 2),
('d5111628-d3a9-11ed-afa1-0242ac120002', 3);

-- Insert default data for ProductBill pivot table
INSERT INTO ProductBill (product_id, bill_id, quantity) VALUES
(1, 1, 2),
(2, 1, 1),
(4, 1, 3),
(5, 1, 2),
(8, 1, 4),
(1, 2, 1),
(3, 2, 1),
(5, 2, 3),
(6, 2, 2),
(7, 2, 1),
(2, 3, 2),
(4, 3, 1),
(6, 3, 1),
(8, 3, 3),
(9, 3, 4),
(3, 4, 3),
(5, 4, 1),
(7, 4, 2),
(9, 4, 2),
(10, 4, 3),
(1, 5, 1),
(2, 5, 2),
(4, 5, 1),
(6, 5, 1),
(8, 5, 2),
(10, 5, 3),
(2, 6, 1),
(3, 6, 2),
(5, 6, 1),
(7, 6, 4),
(9, 6, 1),
(1, 7, 3),
(4, 7, 2),
(6, 7, 1),
(8, 7, 3),
(10, 7, 2),
(1, 8, 1),
(2, 8, 1),
(3, 8, 1),
(5, 8, 2),
(7, 8, 3),
(9, 8, 1);