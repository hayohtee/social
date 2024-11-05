package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hayohtee/social/internal/model"
)

type UsersRepository struct {
	db *sql.DB
}

func (u *UsersRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users(username, email, password)
		VALUES($1, $2, $3)
		RETURNING id, created_at`

	return u.db.QueryRowContext(ctx, query, user.Username, user.Email, user.Password).Scan(
		&user.ID,
		&user.CreatedAt,
	)
}

func (u *UsersRepository) GetByID(ctx context.Context, id int64) (model.User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user model.User
	err := u.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.User{}, ErrNotFound
		default:
			return model.User{}, err
		}
	}
	return user, nil
}
