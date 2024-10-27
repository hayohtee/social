package repository

import (
	"context"
	"database/sql"

	"github.com/hayohtee/social/internal/data"
)

type Repository struct {
	Posts interface {
		Create(context.Context, *data.Post) error
	}
	Users interface {
		Create(context.Context, *data.User) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Posts: &PostsRepository{db: db},
		Users: &UsersRepository{db: db},
	}
}
