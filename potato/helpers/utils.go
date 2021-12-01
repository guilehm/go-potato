package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/guilehm/go-potato/models"

	"github.com/bwmarrin/discordgo"
)

func MakeEmbed(
	url,
	title,
	description string,
	image *discordgo.MessageEmbedImage,
	fields []*discordgo.MessageEmbedField,
	thumbnail *discordgo.MessageEmbedThumbnail,
) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		URL:         url,
		Type:        "",
		Title:       title,
		Description: description,
		Timestamp:   time.Now().Format("2006-01-02 15:04"),
		Color:       3447003,
		Footer: &discordgo.MessageEmbedFooter{
			IconURL:      "",
			Text:         "go potato",
			ProxyIconURL: "",
		},
		Image:     image,
		Thumbnail: thumbnail,
		Video:     nil,
		Provider:  nil,
		Author: &discordgo.MessageEmbedAuthor{
			Name:         "the movie db",
			IconURL:      "https://www.themoviedb.org/assets/2/apple-touch-icon-57ed4b3b0450fd5e9a0c20f34e814b82adaa1085c79bdde2f00ca8787b63d2c4.png",
			URL:          "https://www.themoviedb.org/",
			ProxyIconURL: "",
		},
		Fields: fields,
	}
}

func GetSimpleEmbedForTVShow(tvShow models.TVShowResult) *discordgo.MessageEmbed {
	thumbnail := &discordgo.MessageEmbedThumbnail{
		URL:      "https://www.themoviedb.org/t/p/w300" + tvShow.BackdropPath,
		ProxyURL: "",
		Width:    300,
		Height:   169,
	}

	return &discordgo.MessageEmbed{
		URL: fmt.Sprintf(
			"%v/%v-%v",
			"https://www.themoviedb.org/tv",
			tvShow.ID,
			strings.ReplaceAll(tvShow.Name, " ", "-"),
		),
		Type:        "",
		Title:       tvShow.Name,
		Description: "",
		Timestamp:   "",
		Color:       15418782,
		Footer:      nil,
		Image:       nil,
		Thumbnail:   thumbnail,
		Video:       nil,
		Provider:    nil,
		Author:      nil,
		Fields:      nil,
	}

}

func GetEmbedForTVShow(tvShow models.TVShowResult) *discordgo.MessageEmbed {
	embedImage := &discordgo.MessageEmbedImage{
		URL:      "https://www.themoviedb.org/t/p/w300" + tvShow.BackdropPath,
		ProxyURL: "",
		Width:    300,
		Height:   169,
	}
	thumbnail := &discordgo.MessageEmbedThumbnail{
		URL:      "https://www.themoviedb.org/t/p/w300" + tvShow.PosterPath,
		ProxyURL: "",
		Width:    300,
		Height:   169,
	}

	embedFields := []*discordgo.MessageEmbedField{
		{
			Name:   "Status",
			Value:  tvShow.Status,
			Inline: true,
		},
		{
			Name:   "User Score",
			Value:  fmt.Sprintf("%.0f", tvShow.VoteAverage*10) + "%",
			Inline: true,
		},
		{
			Name:   "No. of Seasons",
			Value:  strconv.Itoa(tvShow.NumberOfSeasons),
			Inline: true,
		},
	}

	if tvShow.Tagline != "" {
		embedFields = append(
			[]*discordgo.MessageEmbedField{
				{
					Name:   "Tagline",
					Value:  tvShow.Tagline,
					Inline: false,
				},
			},
			embedFields...,
		)
	}

	return MakeEmbed(
		fmt.Sprintf(
			"%v/%v-%v",
			"https://www.themoviedb.org/tv",
			tvShow.ID,
			strings.ReplaceAll(tvShow.Name, " ", "-"),
		),
		tvShow.Name,
		tvShow.Overview,
		embedImage,
		embedFields,
		thumbnail,
	)

}

func MakeMovieSearchResultTitles(mr models.MovieSearchResponse) string {
	resultTitles := make([]string, len(mr.Results))
	for index, result := range mr.Results {
		resultTitles[index] = fmt.Sprintf("%v *(%v)*", result.Title, result.ID)
	}
	return strings.Join(resultTitles, "\n")
}

func MakeTVShowSearchResultTitles(sr models.TVSearchResponse) string {
	resultTitles := make([]string, len(sr.Results))
	for index, result := range sr.Results {
		resultTitles[index] = fmt.Sprintf("%v *(%v)*", result.Name, result.ID)
	}
	return strings.Join(resultTitles, "\n")
}
