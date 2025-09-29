package model

type Cart struct {
	UserID UserID
	Items  map[Sku]uint32
}

type ReportCart struct {
	Items      []ItemInCart
	TotalPrice uint32
}

type ItemInCart struct {
	SKU   Sku
	Count uint32
	Name  string
	Price uint32
}
