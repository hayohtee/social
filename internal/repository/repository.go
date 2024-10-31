package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hayohtee/social/internal/data"
)

var (
	ErrNotFound = errors.New("record not found")
)

type Repository struct {
	Posts interface {
		Create(context.Context, *data.Post) error
		GetByID(context.Context, int64) (data.Post, error)
	}

	Users interface {
		Create(context.Context, *data.User) error
	}

	Comments interface {
		GetByPostID(context.Context, int64) ([]data.CommentWithUser, error)
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Posts:    &PostsRepository{db: db},
		Users:    &UsersRepository{db: db},
		Comments: &CommentsRepository{db: db},
	}
}
