package repository

import (
	"context"
	"database/sql"
)

type PostsRepository struct {
	db *sql.DB
}

func (p *PostsRepository) Create(ctx context.Context) error {
	return nil
}
