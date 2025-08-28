package model

type Cart struct {
	UserID UserID
	Items  map[Sku]uint16
}

type ReportCart struct {
	Items      []ItemInСart
	TotalPrice uint32
}

type ItemInСart struct {
	SKU   Sku
	Count uint16
	Name  string
	Price uint32
}
