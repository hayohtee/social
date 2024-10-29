package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hayohtee/social/internal/data"
)

type PostsRepository struct {
	db *sql.DB
}

func (p *PostsRepository) Create(ctx context.Context, post *data.Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	args := []any{post.Content, post.Title, post.UserID, post.Tags}

	return p.db.QueryRowContext(ctx, query, args...).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}

func (p *PostsRepository) GetByID(ctx context.Context, id int64) (data.Post, error) {
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at 
		FROM posts
		WHERE id = $1`

	var post data.Post
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.Tags,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return data.Post{}, ErrNotFound
		default:
			return data.Post{}, err
		}
	}

	return post, nil
}
