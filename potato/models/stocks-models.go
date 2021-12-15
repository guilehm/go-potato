package models

import (
	"encoding/json"
	"time"
)

type StockPriceResponse struct {
	By            string                     `json:"by"`
	ValidKey      bool                       `json:"valid_key"`
	Results       map[string]json.RawMessage `json:"results"`
	Stock         *Stock                     `json:"result"`
	ExecutionTime float64                    `json:"execution_time"`
	FromCache     bool                       `json:"from_cache"`
}

type Stock struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	CompanyName string `json:"company_name"`
	Document    string `json:"document"`
	Description string `json:"description"`
	Website     string `json:"website"`
	Region      string `json:"region"`
	Currency    string `json:"currency"`
	MarketTime  struct {
		Open     string `json:"open"`
		Close    string `json:"close"`
		Timezone int    `json:"timezone"`
	} `json:"market_time"`
	MarketCap     float64   `json:"market_cap"`
	Price         float64   `json:"price"`
	ChangePercent float64   `json:"change_percent"`
	UpdatedAt     time.Time `json:"updated_at"`
	Error         bool      `json:"error"`
}
