package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BaseApiUrl = "https://api.themoviedb.org/3/"
var ErrNotFound = errors.New("not found")

type tmdbService struct {
	ApiKey string
	AccessToken string
}

func (t *tmdbService) Search(text string) error {
	return nil
}

func (t *tmdbService) makeRequest(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v%v", BaseApiUrl, endpoint), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", t.AccessToken))
	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := http.DefaultClient.Do(req)

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	defer resp.Body.Close()

	body, err :=  ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

