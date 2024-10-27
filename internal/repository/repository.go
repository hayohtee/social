package repository

import (
	"context"
	"database/sql"
)

type Repository struct {
	Posts interface {
		Create(context.Context) error
	}
	Users interface {
		Create(context.Context) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Posts: &PostsRepository{db: db},
		Users: &UsersRepository{db: db},
	}
}
