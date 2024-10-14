package postgres

import (
	"database/sql"
	"fmt"
	"io"
	"time"
)

type PostgresRepository interface {
	io.Closer
}

type dbReadWriter struct {
	db *sql.DB
}

func NewPostgres(Host, User, Password, DBName string) (PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, 0000, User, Password, DBName)
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
