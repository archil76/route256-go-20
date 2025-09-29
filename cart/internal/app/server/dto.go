package server

type AddItemRequest struct {
	Count int32 `json:"count"  validate:"min=1"`
}

type ReportCart struct {
	Items      []ItemInСart `json:"items"`
	TotalPrice int32        `json:"total_price"`
}

type ItemInСart struct {
	SKU   int64  `json:"sku"`
	Count int32  `json:"count"`
	Name  string `json:"name"`
	Price int32  `json:"price"`
}

type CheckoutResponse struct {
	OrderID int64 `json:"order_id"`
}
