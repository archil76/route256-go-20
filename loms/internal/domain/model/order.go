package model

type Order struct {
	OrderId int64
	UserID  int64
	Status  Status
	Items   []Item
}
