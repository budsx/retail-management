package postgres

import (
	"context"
	"database/sql"
	reflect "reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/budsx/retail-management/model"
)

func Test_RegisterUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	tests := []struct {
		name    string
		db      *sql.DB
		user    model.User
		mock    func(sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "Successfully register user",
			db:   mockDB,
			user: model.User{
				Username: "testuser",
				Password: "hashed_password",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_users (username, password_hash, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`)).
					WithArgs("testuser", "hashed_password").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "Duplicate username",
			db:   mockDB,
			user: model.User{
				Username: "existing_user",
				Password: "hashed_password",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_users (username, password_hash, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`)).
					WithArgs("existing_user", "hashed_password").
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
		{
			name: "Context cancelled",
			db:   mockDB,
			user: model.User{
				Username: "testuser",
				Password: "hashed_password",
			},
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO mst_users (username, password_hash, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)`)).
					WithArgs("testuser", "hashed_password").
					WillReturnError(context.Canceled)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock(mock)
			}

			rw := &dbReadWriter{
				db: tt.db,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := rw.RegisterUser(ctx, tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbReadWriter.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func Test_GetUserByUsername(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer mockDB.Close()

	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		username string
		mock     func(sqlmock.Sqlmock)
		want     model.User
		wantErr  bool
	}{
		{
			name:     "Successfully retrieve user",
			username: "testuser",
			mock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"user_id", "username", "password_hash", "created_at"}).
					AddRow(1, "testuser", "hashedpassword", fixedTime)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id, username, password_hash, created_at FROM mst_users WHERE username = $1`)).
					WithArgs("testuser").
					WillReturnRows(rows)
			},
			want: model.User{
				UserID:    1,
				Username:  "testuser",
				Password:  "hashedpassword",
				CreatedAt: fixedTime,
			},
			wantErr: false,
		},
		{
			name:     "User not found",
			username: "nonexistent",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id, username, password_hash, created_at FROM mst_users WHERE username = $1`)).
					WithArgs("nonexistent").
					WillReturnError(sql.ErrNoRows)
			},
			want:    model.User{},
			wantErr: true,
		},
		{
			name:     "Database error",
			username: "testuser",
			mock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id, username, password_hash, created_at FROM mst_users WHERE username = $1`)).
					WithArgs("testuser").
					WillReturnError(sql.ErrConnDone)
			},
			want:    model.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mock != nil {
				tt.mock(mock)
			}

			rw := &dbReadWriter{
				db: mockDB,
			}

			got, err := rw.GetUserByUsername(context.Background(), tt.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("dbReadWriter.GetUserByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dbReadWriter.GetUserByUsername() = %v, want %v", got, tt.want)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
