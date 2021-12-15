package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/guilehm/go-potato/potato/models"
)

const BaseStockAPIURL = "https://api.hgbrasil.com/finance/"

type StocksService struct {
	SecretKey string
}

func (s StocksService) makeRequest(endpoint string, queries url.Values) ([]byte, error) {
	u, err := url.Parse(fmt.Sprintf("%v%v", BaseStockAPIURL, endpoint))
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("key", os.Getenv("STOCKS_API_SECRET_KEY"))
	for key, values := range queries {
		for _, v := range values {
			q.Set(key, v)
		}
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil

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
