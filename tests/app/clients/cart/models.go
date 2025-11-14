package cart

type addItemRequest struct {
	Count int64 `json:"count"`
}

type listResponse struct {
	Items []struct {
		SKU   uint32 `json:"sku"`
		Count int64  `json:"count"`
		Name  string `json:"name"`
		Price uint32 `json:"price"`
	} `json:"items"`
	TotalPrice uint32 `json:"total_price"`
}

type checkoutRequest struct {
	User int64 `json:"user_id"`
}

type checkoutResponse struct {
	OrderID int64 `json:"order_id"`
}
