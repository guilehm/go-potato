package services

import (
	"errors"
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
