package model

import "errors"

var (
	ErrSkuIsNotValid       = errors.New("sku should be more than 0")
	ErrCommentDoesntExist  = errors.New("comment doesn't exist")
	ErrCommentIDIsNotValid = errors.New("comment ID should be more than 0")
	ErrUserIDIsNotValid    = errors.New("user ID should be more than 0")
	ErrUserNotAuthor       = errors.New("user is not the author of the comment")
	ErrEditTimeExpired     = errors.New("edit time has expired")
)
