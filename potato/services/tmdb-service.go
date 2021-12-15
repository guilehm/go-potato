package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/guilehm/go-potato/potato/models"
)

const BaseAPIURL = "https://api.themoviedb.org/3/"
const BaseSiteURL = "https://www.themoviedb.org/"

var ErrNotFound = errors.New("not found")

type TMDBService struct {
	AccessToken string
}

func (t *TMDBService) makeSearch(i interface{}, page int, text, endpoint string) error {
	queries := url.Values{
		"query": []string{url.QueryEscape(text)},
		"page":  []string{strconv.Itoa(page)},
	}

	body, err := t.makeRequest(endpoint, queries)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &i); err != nil {
		return err
	}
	return nil

}

func (t *TMDBService) SearchMovies(text string, page int) (models.MovieSearchResponse, error) {
	var response models.MovieSearchResponse
	err := t.makeSearch(&response, page, text, "search/movie")
	return response, err
}

func (t *TMDBService) SearchTvShows(text string, page int) (models.TVSearchResponse, error) {
	var response models.TVSearchResponse
	err := t.makeSearch(&response, page, text, "search/tv")
	return response, err
}

func (t *TMDBService) GetTVShowDetail(id string) (models.TVShowResult, error) {
	var tvShow models.TVShowResult
	q := url.Values{"append_to_response": []string{"credits"}}
	body, err := t.makeRequest("tv/"+id, q)
	if err != nil {
		return tvShow, err
	}

	if err = json.Unmarshal(body, &tvShow); err != nil {
		return tvShow, err
	}
	return tvShow, nil
}

func (t *TMDBService) GetMovieDetail(id string) (models.MovieResult, error) {
	var movie models.MovieResult
	q := url.Values{"append_to_response": []string{"credits"}}
	body, err := t.makeRequest("movie/"+id, q)
	if err != nil {
		return movie, err
	}

	if err = json.Unmarshal(body, &movie); err != nil {
		return movie, err
	}
	return movie, nil
}

func (t *TMDBService) makeRequest(endpoint string, queries url.Values) ([]byte, error) {
	u, err := url.Parse(fmt.Sprintf("%v%v", BaseAPIURL, endpoint))
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
