package loms

type orderItem struct {
	Sku   string `json:"sku"`
	Count int64  `json:"count"`
}

type orderCreateRequest struct {
	UserID int64       `json:"userId,string"`
	Items  []orderItem `json:"items"`
}

type orderCreateResponse struct {
	OrderID int64 `json:"orderId,string"`
}

type orderRequest struct {
	OrderID int64 `json:"orderId,string"`
}

type OrderInfoResponse struct {
	Status string      `json:"status"`
	User   int64       `json:"userId,string"`
	Items  []orderItem `json:"items"`
}

type stocksInfoResponse struct {
	Count uint64 `json:"count"`
}
