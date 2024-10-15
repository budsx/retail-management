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
    user_id INT REFERENCES users(user_id),
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
    product_id INT NOT NULL,               -- Produk yang ada di lokasi ini
    warehouse_id INT NOT NULL,             -- Warehouse atau lokasi spesifik
    stock_quantity INT NOT NULL CHECK (stock_quantity > 0), -- Jumlah stok (tidak boleh nol atau negatif)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (product_id) REFERENCES mst_product(product_id),
    FOREIGN KEY (warehouse_id) REFERENCES mst_warehouse(warehouse_id)
);

-- Inventory
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

-- User Register, Login and Validate token
-- User CRUD Product √
-- User CRUD Location berdasarkan warehouse yang dia punya saja
-- User dapat update inventory yang dia punya
    -- Update trx_stock
    -- Create trx_inventory
    -- Read total stock all location / single location
    -- Read transaksi dia sendiri