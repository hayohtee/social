package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hayohtee/social/internal/data"
)

type UsersRepository struct {
	db *sql.DB
}

func (u *UsersRepository) Create(ctx context.Context, user *data.User) error {
	query := `
		INSERT INTO users(username, email, password)
		VALUES($1, $2, $3)
		RETURNING id, created_at`

	return u.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(
		&user.ID,
		&user.CreatedAt,
	)
}

func (u *UsersRepository) GetByID(ctx context.Context, id int64) (data.User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user data.User
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return data.User{}, ErrNotFound
		default:
			return data.User{}, err
		}
	}
	return user, nil
}
