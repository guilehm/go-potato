package services

import (
	"encoding/json"
	"fmt"

	"github.com/guilehm/go-potato/potato/models"
)

const BaseStockAPIURL = "https://api.hgbrasil.com/finance/"

type StocksService struct {
	SecretKey string
}

func (s StocksService) unmarshallStockPriceResponse(body []byte) (*models.StockPriceResponse, error) {
	r := &models.StockPriceResponse{}
	err := json.Unmarshal(body, r)
	if err != nil {
		fmt.Println("Could not unmarshall stock data")
		return r, err
	}

	// Only the first result matters
	for _, v := range r.Results {
		stock := &models.Stock{}
		if err := json.Unmarshal(v, stock); err != nil {
			fmt.Println("Could not unmarshall stock data")
			return r, err
		}
		r.Stock = stock
		break
	}
	return r, nil
}
