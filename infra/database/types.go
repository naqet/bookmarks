package database

import "time"

type User struct {
	ID        string
	Username  string
	Password  string
	CreatedAt time.Time
}

type Bookmark struct {
	ID          string
	Title       string
	Url         string
	Tags        string
	Description string
	Read        bool
	CreatedAt   time.Time
}
