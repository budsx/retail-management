package postgres

import (
	"context"

	"github.com/budsx/retail-management/model"
)

func (rw *dbReadWriter) RegisterUser(ctx context.Context, user model.User) error {
	query := `INSERT INTO mst_users (username, password_hash, created_at) 
              VALUES ($1, $2, CURRENT_TIMESTAMP)`

	_, err := rw.db.ExecContext(ctx, query, user.Username, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (rw *dbReadWriter) GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User

	query := `SELECT user_id, username, password_hash, created_at 
              FROM mst_users 
              WHERE username = $1`

	err := rw.db.QueryRowContext(ctx, query, username).Scan(
		&user.UserID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
