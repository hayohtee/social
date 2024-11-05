package model

import "time"

type Feed struct {
	PostID        int64
	Title         string
	UserID        int64
	Tags          []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Content       string
	Username      string
	CommentsCount int
}
