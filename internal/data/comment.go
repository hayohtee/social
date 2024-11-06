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
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    int64     `json:"user_id"`
	Content   string    `json:"content"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}
