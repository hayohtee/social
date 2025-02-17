package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/hayohtee/social/internal/data"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrEditConflict      = errors.New("edit conflict")
	ErrDuplicateKey      = errors.New("resource already exist")
	ErrDuplicateEmail    = errors.New("a user with this email already exist")
	ErrDuplicateUsername = errors.New("a user with this username already exist")
)

const queryTimeoutDuration = 5 * time.Second

type Repository struct {
	Posts interface {
		Create(context.Context, *data.Post) error
		GetByID(context.Context, int64) (data.Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *data.Post) error
		GetUserFeeds(context.Context, int64, data.Filters) ([]data.Feed, error)
	}

	Users interface {
		Create(context.Context, *data.User, *sql.Tx) error
		CreateAndInvite(context.Context, *data.User, []byte, time.Duration) error
		Activate(context.Context, string) error
		GetByID(context.Context, int64) (data.User, error)
	}

	Comments interface {
		Create(context.Context, *data.Comment) error
		GetByPostID(context.Context, int64) ([]data.CommentWithUser, error)
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
