package main

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"strings"
)

//go:embed stock-data.json
var stockData string

type StockData struct {
	Sku        int64  `json:"sku"`
	TotalCount uint32 `json:"total_count"`
	Reserved   uint32 `json:"reserved"`
}

func loadSource() ([]StockData, error) {
	jsonData := []byte(stockData)
	var stockData []StockData
	if err := json.Unmarshal(jsonData, &stockData); err != nil {
		return nil, err
	}

	return stockData, nil
}

func main() {
	stockData, err := loadSource()
	if err != nil {
		return
	}

	sb := strings.Builder{}
	sb.WriteString("INSERT INTO stocks (id, total_count, reserved)")
	sb.WriteString("\n")
	sb.WriteString("VALUES ")
	sb.WriteString("\n")
	for _, stock := range stockData {
		line := fmt.Sprintf("(%d, %d, %d),", stock.Sku, stock.TotalCount, stock.Reserved)
		if err != nil {
			return
		}
		sb.WriteString(line)
		sb.WriteString("\n")
	}
	sb.WriteString(";")

	fmt.Print(sb.String())
}
