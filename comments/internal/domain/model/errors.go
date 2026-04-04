package model

import "errors"

var (
	ErrStockDoesntExist = errors.New("stock doesn't exist")
	ErrSkuIsNotValid    = errors.New("sku should be more than 0")
	ErrOutOfStock       = errors.New("available amount of stock isn't enough ")

	ErrCommentDoesntExist  = errors.New("order doesn't exist")
	ErrCommentIDIsNotValid = errors.New("comment ID should be more than 0")

	ErrOrderIDIsNotValid  = errors.New("order ID should be more than 0")
	ErrUserIDIsNotValid   = errors.New("user ID should be more than 0")
	ErrInvalidOrderStatus = errors.New("order status should be PAYED")
)
