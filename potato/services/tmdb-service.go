package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/guilehm/go-potato/models"
)

const BaseApiUrl = "https://api.themoviedb.org/3/"

var ErrNotFound = errors.New("not found")

type TMDBService struct {
	AccessToken string
}

func (t *TMDBService) SearchMovies(text string, page int) (models.MovieSearchResponse, error) {
	var response models.MovieSearchResponse
	queries := url.Values{
		"query": []string{url.QueryEscape(text)},
		"page":  []string{strconv.Itoa(page)},
	}

	body, err := t.makeRequest("search/movie", queries)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

func (t *TMDBService) SearchTvShows(text string, page int) (models.TVSearchResponse, error) {
	var response models.TVSearchResponse

	queries := url.Values{
		"query": []string{url.QueryEscape(text)},
		"page":  []string{strconv.Itoa(page)},
	}

	body, err := t.makeRequest("search/tv", queries)
	if err != nil {
		return response, err
	}

	if err = json.Unmarshal(body, &response); err != nil {
		return response, err
	}
	return response, nil
}

func (t *TMDBService) GetTVShowDetail(id string) (models.TVShow, error) {
	var tvShow models.TVShow
	body, err := t.makeRequest("tv/"+id, nil)
	if err != nil {
		return tvShow, err
	}

	if err = json.Unmarshal(body, &tvShow); err != nil {
		return tvShow, err
	}
	return tvShow, nil
}

func (t *TMDBService) makeRequest(endpoint string, queries url.Values) ([]byte, error) {
	u, err := url.Parse(fmt.Sprintf("%v%v", BaseApiUrl, endpoint))
	if err != nil {
		return nil, err
	}

	q := u.Query()
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
