package inmemoryrepository

import (
	"encoding/json"
	"errors"
	"os"
	"route256/loms/internal/domain/model"
	"sync"
)

var (
	ErrStockDoesntExist = errors.New("stock doesn't exist")
	ErrSkuIsNotValid    = errors.New("sku should be more than 0")
	ErrShortOfStock     = errors.New("available amount of stock isn't enough ")
)

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

func loadSource(filename string) ([]StockData, error) {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var stockData []StockData
	if err := json.NewDecoder(f).Decode(&stockData); err != nil {
		return nil, err
	}

	return stockData, nil
}

func NewStockInMemoryRepository(source string, capacity int) *Repository {
	stockData, err := loadSource(source)
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
