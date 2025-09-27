package model

type Cart struct {
	UserID UserID
	Items  map[Sku]uint32
}

type ReportCart struct {
	Items      []ItemInСart
	TotalPrice uint32
}

type ItemInСart struct {
	SKU   Sku
	Count uint32
	Name  string
	Price uint32
}
