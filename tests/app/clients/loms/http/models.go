package http

type OrderItem struct {
	Sku   uint32 `json:"sku"`
	Count uint32 `json:"count"`
}

type OrderCreateRequest struct {
	User  int64       `json:"user"`
	Items []OrderItem `json:"items"`
}

type OrderCreateResponse struct {
	OrderId int64 `json:"orderId,string"`
}

type OrderInfoRequest struct {
	OrderID int64 `json:"orderId"`
}

type OrderInfoResponse struct {
	Status string      `json:"status"`
	User   int64       `json:"user,string"`
	Items  []OrderItem `json:"items"`
}

type OrderPayRequest struct {
	OrderID int64 `json:"orderId"`
}

type OrderPayResponse struct {
}

type OrderCancelRequest struct {
	OrderID int64 `json:"orderId"`
}

type OrderCancelResponse struct {
}

type StocksInfoRequest struct {
	Sku uint32 `json:"sku"`
}

type StocksInfoResponse struct {
	Count uint64 `json:"count,string"`
}
