package data

import "time"

type Post struct {
	ID        int64
	Content   string
	Title     string
	UserID    int64
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
