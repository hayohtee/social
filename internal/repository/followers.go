package repository

import (
	"context"
	"database/sql"
)

type FollowersRepository struct {
	db *sql.DB
}

func (f *FollowersRepository) UnFollow(ctx context.Context, userID int64, followerID int64) error {
	return nil
}

func (f *FollowersRepository) Follow(ctx context.Context, userID int64, followerID int64) error {
	return nil
}
