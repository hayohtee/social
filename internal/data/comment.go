package data

import "time"

type Comment struct {
	ID        int64
	PostID    int64
	UserID    int64
	Content   string
	CreatedAt time.Time
}

type CommentWithUser struct {
	ID        int64
	PostID    int64
	UserID    int64
	Content   string
	UserName  string
	CreatedAt time.Time
}
