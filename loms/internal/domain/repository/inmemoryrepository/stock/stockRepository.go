package inmemoryrepository

import (
	_ "embed"
	"encoding/json"
	"route256/loms/internal/domain/model"
	"sync"
)

//go:embed stock-data.json
var stockData string

type StockData struct {
	Sku        int64  `json:"sku"`
	TotalCount uint32 `json:"total_count"`
	Reserved   uint32 `json:"reserved"`
}

type Storage = map[int64]model.Stock

type Repository struct {
	storage Storage
	mu      sync.RWMutex
}

func loadSource() ([]StockData, error) {
	jsonData := []byte(stockData)
	var stockData []StockData
	if err := json.Unmarshal(jsonData, &stockData); err != nil {
		return nil, err
	}

	return stockData, nil
}

func NewStockInMemoryRepository(capacity int) *Repository {
	stockData, err := loadSource()
	if err != nil {
		return nil
	}

	repository := &Repository{storage: make(Storage, max(capacity, len(stockData)))}

	for _, stockData := range stockData {
		repository.storage[stockData.Sku] = model.Stock{
			Sku:        stockData.Sku,
			TotalCount: stockData.TotalCount,
			Reserved:   stockData.Reserved,
		}
	}

	return repository
}
