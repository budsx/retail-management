package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/budsx/retail-management/model"
)

func (rw *dbReadWriter) ReadProductByID(ctx context.Context, req int64) (model.Product, error) {
	selectProductByID := `SELECT product_id, product_name, description, price, sku, created_at, updated_at 
	FROM mst_product 
	WHERE product_id = $1`

	var product model.Product
	err := rw.db.QueryRowContext(ctx, selectProductByID, req).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.Description,
		&product.Price,
		&product.SKU,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return product, fmt.Errorf("product with id %d not found", req)
		}
		return product, err
	}

	return product, nil
}

func (rw *dbReadWriter) ReadProductsWithPagination(ctx context.Context, limit int32, offset int32) ([]model.Product, error) {
	selectProductsWithPagination := `SELECT product_id, product_name, description, price, sku, created_at, updated_at 
		FROM mst_product
		LIMIT $1 OFFSET $2`

	rows, err := rw.db.QueryContext(ctx, selectProductsWithPagination, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product

	for rows.Next() {
		var product model.Product
		if err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.Description,
			&product.Price,
			&product.SKU,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}


func (rw *dbReadWriter) UpdateProductByID(ctx context.Context, product model.Product) error {
	updateProduct := `UPDATE mst_product 
		SET product_name = $1, description = $2, price = $3, sku = $4, updated_at = CURRENT_TIMESTAMP 
		WHERE product_id = $5`

	result, err := rw.db.ExecContext(ctx, updateProduct,
		product.ProductName,
		product.Description,
		product.Price,
		product.SKU,
		product.ProductID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with id %d not found", product.ProductID)
	}

	return nil
}

func (rw *dbReadWriter) WriteProduct(ctx context.Context, product model.Product) error {
	insertProduct := `INSERT INTO mst_product (product_name, description, price, sku, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	_, err := rw.db.ExecContext(ctx, insertProduct,
		product.ProductName,
		product.Description,
		product.Price,
		product.SKU,
	)

	if err != nil {
		return err
	}

	return nil
}
