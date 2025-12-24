package comments

import "time"

type addCommentRequest struct {
	UserID  int64  `json:"user_id"`
	Sku     int64  `json:"sku"`
	Comment string `json:"comment"`
}

type addCommentResponse struct {
	CommentID int64 `json:"id,string"`
}

type editCommentRequest struct {
	UserID     int64  `json:"user_id"`
	CommentID  int64  `json:"comment_id"`
	NewComment string `json:"new_comment"`
}

type getCommentListBySKURequest struct {
	Sku int64 `json:"sku"`
}

type getCommentListByUserRequest struct {
	UserID int64 `json:"user_id"`
}

type commentGetByIDRequest struct {
	ID int64 `json:"id"`
}

type comments struct {
	Comments []comment `json:"comments"`
}

type comment struct {
	ID        int64     `json:"id,string"`
	UserID    int64     `json:"userId,string"`
	SKU       int64     `json:"sku,string"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"createdAt"`
}
