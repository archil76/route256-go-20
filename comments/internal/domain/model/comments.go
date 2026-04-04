package model

import (
	"time"
)

type Comment struct {
	ID        int64
	UserID    int64
	Sku       int64
	CreatedAt time.Time
	Comment   string
}

type CommentsList struct {
	Comments []Comment
}
