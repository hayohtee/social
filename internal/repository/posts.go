package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/hayohtee/social/internal/model"
	"github.com/lib/pq"
)

type PostsRepository struct {
	db *sql.DB
}

func (p *PostsRepository) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	args := []any{post.Content, post.Title, post.UserID, pq.Array(post.Tags)}

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	return p.db.QueryRowContext(ctx, query, args...).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}

func (p *PostsRepository) GetByID(ctx context.Context, id int64) (model.Post, error) {
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var post model.Post
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return model.Post{}, ErrNotFound
		default:
			return model.Post{}, err
		}
	}

	return post, nil
}

func (p *PostsRepository) Delete(ctx context.Context, postID int64) error {
	query := `
		DELETE FROM posts
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	row, err := p.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (p *PostsRepository) Update(ctx context.Context, post *model.Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, updated_at = NOW(), version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING updated_at, version`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	args := []any{post.Title, post.Content, post.ID, post.Version}
	err := p.db.QueryRowContext(ctx, query, args...).Scan(
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (p *PostsRepository) GetUserFeeds(ctx context.Context, userID int64) ([]model.Feed, error) {
	query := `
		SELECT 
			posts.id AS post_id,
			posts.title,
			posts.content,
			posts.user_id AS author_id,
			users.username AS author_username,
			posts.tags,
			posts.created_at,
			posts.updated_at,
			COUNT(comments.id) AS comment_count
		FROM 
			posts
		JOIN 
			followers ON posts.user_id = followers.user_id
		JOIN 
			users ON posts.user_id = users.id
		LEFT JOIN 
			comments ON posts.id = comments.post_id
		WHERE 
			followers.follower_id = $1
		GROUP BY 
			posts.id, 
			posts.title, 
			posts.content, 
			posts.user_id, 
			users.username, 
			posts.tags, 
			posts.created_at
		ORDER BY 
			posts.created_at DESC`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var feeds []model.Feed

	rows, err := p.db.QueryContext(ctx, query, userID)
	if err != nil {
		return []model.Feed{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var feed model.Feed
		err := rows.Scan(
			&feed.PostID,
			&feed.Title,
			&feed.Content,
			&feed.AuthorID,
			&feed.AuthorUsername,
			pq.Array(&feed.Tags),
			&feed.CreatedAt,
			&feed.UpdatedAt,
			&feed.CommentsCount,
		)

		if err != nil {
			return []model.Feed{}, err
		}

		feeds = append(feeds, feed)
	}

	if rows.Err() != nil {
		return []model.Feed{}, err
	}

	return feeds, nil
}
