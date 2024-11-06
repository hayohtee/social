package data

import "time"

type Feed struct {
	PostID         int64     `json:"post_id"`
	Title          string    `json:"title"`
	AuthorID       int64     `json:"author_id"`
	Tags           []string  `json:"tags"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Content        string    `json:"content"`
	AuthorUsername string    `json:"author_username"`
	CommentsCount  int       `json:"comments_count"`
}
