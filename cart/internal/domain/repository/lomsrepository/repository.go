package lomsrepository

import (
	"errors"
)

var (
	ErrSkuNotFoundInStock = errors.New("sku not found in stock")
	ErrOrderNotFound      = errors.New("product not found")
)

type LomsService struct {
	address string
}

func NewLomsService(address string) *LomsService {
	return &LomsService{
		address: address,
	}
}
