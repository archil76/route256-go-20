package domain

import "time"

type Comment struct {
	ID        int64
	UserID    int64
	SKU       int64
	Comment   string
	CreatedAt time.Time
}
