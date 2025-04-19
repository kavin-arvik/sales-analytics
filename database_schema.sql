CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE customers (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    address TEXT
);

CREATE TABLE products (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    customer_id UUID REFERENCES customers(id),
    date_of_sale DATE,
    region TEXT,
    payment_method TEXT
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id UUID REFERENCES orders(id),
    product_id UUID REFERENCES products(id),
    quantity INT,
    unit_price NUMERIC,
    discount NUMERIC,
    shipping_cost NUMERIC
);


  CREATE TABLE refresh_logs (
    id SERIAL PRIMARY KEY,
    refresh_type VARCHAR(20),
    status VARCHAR(10),      
    started_at TIMESTAMP DEFAULT now(),
    ended_at TIMESTAMP,
    error_message TEXT
);
