package repository

import (
	"context"
	"database/sql"

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
