package model

type Sku = int64

type Item struct {
	Sku   Sku
	Count uint32
}
