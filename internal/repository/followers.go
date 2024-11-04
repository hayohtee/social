package repository

import (
	"context"
	"database/sql"
)

type FollowersRepository struct {
	db *sql.DB
}

func (f *FollowersRepository) UnFollow(ctx context.Context, userID, followerID int64) error {
	query := `
		DELETE FROM followers
		WHERE user_id = $1 AND follower_id = $2`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, userID, followerID)

	return err
}

func (f *FollowersRepository) Follow(ctx context.Context, userID, followerID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, userID, followerID)
	return err
}
