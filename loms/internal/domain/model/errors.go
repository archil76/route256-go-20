package model

import "errors"

var (
	ErrStockDoesntExist = errors.New("stock doesn't exist")
	ErrSkuIsNotValid    = errors.New("sku should be more than 0")
	ErrShortOfStock     = errors.New("available amount of stock isn't enough ")

	ErrOrderDoesntExist = errors.New("order doesn't exist")
	ErrUserIDIsNotValid = errors.New("UserID should be more than 0")

	ErrOrderIDIsNotValid  = errors.New("order ID should be more than 0")
	ErrSkuIDIsNotValid    = errors.New("sku should be more than 0")
	ErrInvalidOrderStatus = errors.New("order status should be PAYED")
)
