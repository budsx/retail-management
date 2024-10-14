package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"time"

	"github.com/budsx/retail-management/model"
	_ "github.com/lib/pq"
)

type PostgresRepository interface {
	ReadProductByID(context.Context, int64) (model.Product, error)
	ReadProductsWithPagination(context.Context, int32, int32) ([]model.Product, error)
	UpdateProductByID(context.Context, model.Product) error
	WriteProduct(context.Context, model.Product) error

	// User
	RegisterUser(context.Context, model.User) error
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	io.Closer
}

type dbReadWriter struct {
	db *sql.DB
}

func NewPostgres(Host, User, Password, DBName string) (PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, 5432, User, Password, DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &dbReadWriter{db}, nil
}

func (rw *dbReadWriter) Close() error {
	if rw.db != nil {
		if err := rw.db.Close(); err != nil {
			return err
		}
		rw.db = nil
	}
	return nil
}
