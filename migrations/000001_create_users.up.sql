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

INSERT INTO public.mst_users
(username, password_hash, created_at)
VALUES('user1', '$2a$10$FXpdvb6bvpGhcvxlGP9/COVBti0le3SH7M5yonYKLpIB3eNZxAv5O', '2024-10-27 14:28:11.380');
INSERT INTO public.mst_users
(username, password_hash, created_at)
VALUES('user2', '$2a$10$DqIz0X2VORWzo66o1mGeTeTZS.THAvv/tjzPOx2Ev3K83DVidaH2u', '2024-10-27 14:28:11.435');
INSERT INTO public.mst_users
(username, password_hash, created_at)
VALUES('user3', '$2a$10$mqZNc3wkR5chYSVR8Zs1VOW0VuGTkFaW0Tm0PzRTIFEnIAQ5VTuX2', '2024-10-27 14:28:11.488');

INSERT INTO mst_product (product_name, description, price, sku) VALUES
('Kopi Arabika', 'Kopi Arabika premium dari Aceh.', 50000.00, 'KOP-AR-001'),
('Teh Hijau', 'Teh hijau organik tanpa gula.', 30000.00, 'TEH-HJ-002'),
('Gula Pasir', 'Gula pasir murni.', 15000.00, 'GUL-PAS-003');

INSERT INTO mst_warehouse (warehouse_name, user_id) VALUES
('Gudang Utama', 1),
('Gudang Cabang Jakarta', 2),
('Gudang Cabang Surabaya', 3);

INSERT INTO mst_location (location_name, warehouse_id) VALUES
('Rak A1', 1),
('Rak B2', 1),
('Rak C3', 2),
('Rak D4', 3);

INSERT INTO mst_stock (product_id, warehouse_id, stock_quantity) VALUES
(1, 1, 100),
(2, 1, 200),
(3, 2, 150),
(1, 3, 80),
(2, 3, 120);

INSERT INTO trx_stock (product_id, warehouse_id, transaction_type, quantity, created_by) VALUES
(1, 1, 'IN', 50, 1),
(2, 1, 'OUT', 10, 2),
(3, 1, 'IN', 20, 3),
(1, 2, 'OUT', 15, 3),
(3, 2, 'IN', 25, 1);

COMMIT;
