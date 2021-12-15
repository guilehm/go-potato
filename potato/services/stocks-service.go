package services

import (
	"encoding/json"
	"fmt"

	"github.com/guilehm/go-potato/potato/models"
)

func UnmarshallStockPriceResponse(r []byte) (*models.StockPriceResponse, error) {
	s := &models.StockPriceResponse{}
	err := json.Unmarshal(r, s)
	if err != nil {
		fmt.Println("Could not unmarshall stock data")
		return s, err
	}

	// Only the first result matters
	for _, v := range s.Results {
		stock := &models.Stock{}
		if err := json.Unmarshal(v, stock); err != nil {
			fmt.Println("Could not unmarshall stock data")
			return s, err
		}
		s.Stock = stock
		break
	}
	return s, nil
}
