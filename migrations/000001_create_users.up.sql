BEGIN;

-- Users
CREATE TABLE mst_users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Products
CREATE TABLE mst_product (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    sku VARCHAR(100) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Warehouse
CREATE TABLE mst_warehouse (
    warehouse_id SERIAL PRIMARY KEY,
    warehouse_name VARCHAR(255) NOT NULL,
    user_id INT REFERENCES mst_users(user_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Location
CREATE TABLE mst_location (
    location_id SERIAL PRIMARY KEY,
    location_name VARCHAR(255) NOT NULL,
    warehouse_id INT REFERENCES mst_warehouse(warehouse_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Stock Based on Location
CREATE TABLE mst_stock (
    stock_id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,               -- Product at this location
    warehouse_id INT NOT NULL,             -- Warehouse or specific location
    stock_quantity INT NOT NULL CHECK (stock_quantity > 0), -- Stock quantity (cannot be zero or negative)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES mst_product(product_id),
    FOREIGN KEY (warehouse_id) REFERENCES mst_warehouse(warehouse_id)
);

-- Stock Transaction
CREATE TABLE trx_stock (
    transaction_id SERIAL PRIMARY KEY,
    product_id INT NOT NULL,
    warehouse_id INT NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    quantity INT NOT NULL,
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT,
    FOREIGN KEY (product_id) REFERENCES mst_product(product_id),
    FOREIGN KEY (warehouse_id) REFERENCES mst_warehouse(warehouse_id)
);

-- Create trigger for updating `updated_at` on `mst_product`
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
   NEW.updated_at = NOW();
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for auto-updating `updated_at` on each update for `mst_product` and `mst_stock`
CREATE TRIGGER update_mst_product_updated_at
BEFORE UPDATE ON mst_product
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_mst_stock_updated_at
BEFORE UPDATE ON mst_stock
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

COMMIT;
