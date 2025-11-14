package domain

type CartItem struct {
	SKU   uint32
	Count int64
	Name  string
	Price uint32
}

type Cart struct {
	Items      []CartItem
	TotalPrice uint32
}
