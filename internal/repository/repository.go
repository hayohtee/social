package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hayohtee/social/internal/model"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
	ErrDuplicateKey = errors.New("resource already exist")

	QueryTimeoutDuration = 5 * time.Second
)

type Repository struct {
	Posts interface {
		Create(context.Context, *model.Post) error
		GetByID(context.Context, int64) (model.Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *model.Post) error
		GetUserFeeds(context.Context, int64) ([]model.Feed, error)
	}

	Users interface {
		Create(context.Context, *model.User) error
		GetByID(context.Context, int64) (model.User, error)
	}

	Comments interface {
		Create(context.Context, *model.Comment) error
		GetByPostID(context.Context, int64) ([]model.CommentWithUser, error)
	}

	Followers interface {
		Follow(context.Context, int64, int64) error
		UnFollow(context.Context, int64, int64) error
	}
}

func NewRepository(db *sql.DB) Repository {
	return Repository{
		Posts:     &PostsRepository{db: db},
		Users:     &UsersRepository{db: db},
		Comments:  &CommentsRepository{db: db},
		Followers: &FollowersRepository{db: db},
	}
}
