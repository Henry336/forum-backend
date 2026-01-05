package models

import "time"

type Topic struct {
	Id   int
	Name string
}

type Post struct {
	Id          int
	Title       string
	Description string
	CreatedAt   time.Time
	TopicId     int
	Username    string // Added from our Auth update
}

type Comment struct {
	Id        int
	Content   string
	CreatedAt time.Time
	PostId    int
	Username  string // Added from our Auth update
}

type User struct {
	Username string
}
