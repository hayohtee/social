package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hayohtee/social/internal/data"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")

	QueryTimeoutDuration = 5 * time.Second
)

type Repository struct {
	Posts interface {
		Create(context.Context, *data.Post) error
		GetByID(context.Context, int64) (data.Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *data.Post) error
	}

	Users interface {
		Create(context.Context, *data.User) error
		GetByID(context.Context, int64) (data.User, error)
	}

	Comments interface {
		Create(context.Context, *data.Comment) error
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
