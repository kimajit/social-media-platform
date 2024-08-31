package models

import (
	"time"
)

// Post represents a social media post.
type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content" validate:"required,max=500"`
	Comments  []Comment `json:"comments"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	CreatedAt time.Time `json:"created_at"`
}

// Comment represents a comment on a post.
type Comment struct {
	ID        string    `json:"id"`
	Content   string    `json:"content" validate:"required,max=300"`
	CreatedAt time.Time `json:"created_at"`
}
