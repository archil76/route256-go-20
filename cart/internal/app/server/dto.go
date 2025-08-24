package server

import (
	model "route256/cart/internal/domain/model"
)

type AddItemRequest struct {
	Count uint16 `json:"count"  validate:"min=1"`
}

type ReportCart struct {
	UserID     model.UserID             `json:"user_id"`
	Items      map[model.Sku]ItemInСart `json:"items"`
	TotalPrice int32                    `json:"total_price"`
}

type ItemInСart struct {
	SKU   model.Sku `json:"sku_id"`
	Count uint16    `json:"count"`
	Name  string    `json:"name"`
	Price int32     `json:"price"`
}
