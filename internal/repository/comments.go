package repository

import (
	"context"
	"database/sql"

	"github.com/hayohtee/social/internal/data"
)

type CommentsRepository struct {
	db *sql.DB
}

func (c *CommentsRepository) GetByPostID(ctx context.Context, postID int64) ([]data.CommentWithUser, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username
		FROM comments c
		JOIN users ON users.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC`

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []data.CommentWithUser
	for rows.Next() {
		var comment data.CommentWithUser
		err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UserName,
		)

		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return comments, nil
}

func (c *CommentsRepository) Create(ctx context.Context, comment *data.Comment) error {
	query := `
		INSERT INTO comments(post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	args := []any{comment.PostID, comment.UserID, comment.Content}

	ctx, cancel := context.WithTimeout(ctx, queryTimeoutDuration)
	defer cancel()

	return c.db.QueryRowContext(ctx, query, args...).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
}
