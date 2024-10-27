-- Dummy
INSERT INTO public.mst_users
(user_id, username, password_hash, created_at)
VALUES(1, 'user1', '$2a$10$FXpdvb6bvpGhcvxlGP9/COVBti0le3SH7M5yonYKLpIB3eNZxAv5O', '2024-10-27 14:28:11.380');
INSERT INTO public.mst_users
(user_id, username, password_hash, created_at)
VALUES(2, 'user2', '$2a$10$DqIz0X2VORWzo66o1mGeTeTZS.THAvv/tjzPOx2Ev3K83DVidaH2u', '2024-10-27 14:28:11.435');
INSERT INTO public.mst_users
(user_id, username, password_hash, created_at)
VALUES(3, 'user3', '$2a$10$mqZNc3wkR5chYSVR8Zs1VOW0VuGTkFaW0Tm0PzRTIFEnIAQ5VTuX2', '2024-10-27 14:28:11.488');

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
(4, 2, 'IN', 100, 1),
(5, 2, 'OUT', 30, 2),
(1, 2, 'OUT', 15, 3),
(3, 2, 'IN', 25, 1);