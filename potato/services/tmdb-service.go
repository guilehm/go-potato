package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/guilehm/go-potato/models"
)

const BaseApiUrl = "https://api.themoviedb.org/3/"

var ErrNotFound = errors.New("not found")

type TMDBService struct {
	AccessToken string
}

func (t *TMDBService) SearchMovies(text string) (models.MovieSearchResponse, error) {
	var response models.MovieSearchResponse
	body, err := t.makeRequest("search/movie?query=" + text)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

func (t *TMDBService) SearchTvShows(text string) (models.TVSearchResponse, error) {
	var response models.TVSearchResponse
	body, err := t.makeRequest("search/tv?query=" + text)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

func (t *TMDBService) makeRequest(endpoint string) ([]byte, error) {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
