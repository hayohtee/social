package repository

import (
	"context"
	"database/sql"
)

type UsersRepository struct {
	db *sql.DB
}

func (u *UsersRepository) Create(ctx context.Context) error {
	return nil
}
