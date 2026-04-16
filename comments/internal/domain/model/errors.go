package model

import "errors"

var (
	ErrSkuIsNotValid       = errors.New("sku should be more than 0")
	ErrCommentDoesntExist  = errors.New("comment doesn't exist")
	ErrCommentIDIsNotValid = errors.New("comment ID should be more than 0")
	ErrUserIDIsNotValid    = errors.New("user ID should be more than 0")
)
