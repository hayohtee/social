package model

import (
	"time"

	"github.com/hayohtee/social/internal/validator"
)

type Post struct {
	ID        int64
	Content   string
	Title     string
	UserID    int64
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
	Version   int
}

func ValidatePost(v *validator.Validator, post Post) {
	v.Check(post.Title != "", "title", "must be provided")
	v.Check(len(post.Title) <= 100, "title", "must not be more than 100 bytes long")

	v.Check(post.Content != "", "content", "must be provided")
	v.Check(len(post.Content) <= 1000, "content", "must be more than 1000 bytes long")

	v.Check(post.UserID > 0, "user_id", "must be a positive number")
}
