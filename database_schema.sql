
-- Customer Table
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    customer_id VARCHAR(50) UNIQUE NOT NULL,
    name TEXT,
    email TEXT,
    address TEXT,
    created_date datetime,
    updated_date TIMESTAMP
);

-- Product Table
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    product_id VARCHAR(50) UNIQUE NOT NULL,
    name TEXT,
    category TEXT,
    created_date datetime,
    updated_date datetime
);

-- Orders Table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(100) UNIQUE NOT NULL,
    customer_id VARCHAR(50),
    region TEXT,
    date_of_sale DATE,
    payment_method TEXT,
    shipping_cost DECIMAL(10,2),
    created_date datetime,
    updated_date datetime,
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
);

-- Order Items Table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id VARCHAR(100),
    product_id VARCHAR(50),
    quantity_sold INT,
    unit_price DECIMAL(10,2),
    discount DECIMAL(5,2),
   	created_date datetime,
    updated_date datetime,
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);

