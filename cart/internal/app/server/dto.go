package server

type AddItemRequest struct {
	Count uint32 `json:"count"  validate:"min=1"`
}

type ReportCart struct {
	Items      []ItemInСart `json:"items"`
	TotalPrice uint32       `json:"total_price"`
}

type ItemInСart struct {
	SKU   int64  `json:"sku"`
	Count uint32 `json:"count"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
