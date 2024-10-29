package data

import "time"

type User struct {
	ID        int64
	Username  string
	Email     string
	Password  string
	CreatedAt time.Time
}
