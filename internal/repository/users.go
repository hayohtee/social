package repository

import (
	"context"
	"database/sql"

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
