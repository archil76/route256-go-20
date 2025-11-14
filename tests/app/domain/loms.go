package domain

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "new"
	OrderStatusAwaitingPayment OrderStatus = "awaiting payment"
	OrderStatusFailed          OrderStatus = "failed"
	OrderStatusPaid            OrderStatus = "paid"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

type OrderItem struct {
	Sku   int32
	Count int64
}

type Order struct {
	Status OrderStatus
	User   int64
	Items  []OrderItem
}

type OrderEvent struct {
	OrderID int64       `json:"order_id"`
	Status  OrderStatus `json:"status"`
}
