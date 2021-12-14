package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/guilehm/go-potato/potato/services"

	"github.com/guilehm/go-potato/potato/models"

	"github.com/bwmarrin/discordgo"
)

const emojiLength = 3

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
			IconURL:      services.BaseSiteURL + "assets/2/apple-touch-icon-57ed4b3b0450fd5e9a0c20f34e814b82adaa1085c79bdde2f00ca8787b63d2c4.png",
			URL:          services.BaseSiteURL,
			ProxyIconURL: "",
		},
		Fields: fields,
	}
}

func GetSimpleEmbedForTVShow(tvShow models.TVShowResult) *discordgo.MessageEmbed {
	thumbnail := &discordgo.MessageEmbedThumbnail{
		URL:      services.BaseSiteURL + "t/p/w300" + tvShow.BackdropPath,
		ProxyURL: "",
		Width:    300,
		Height:   169,
	}

	return &discordgo.MessageEmbed{
		URL: fmt.Sprintf(
			"%v/%v-%v",
			services.BaseSiteURL+"tv",
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
		URL:      services.BaseSiteURL + "t/p/w300" + tvShow.BackdropPath,
		ProxyURL: "",
		Width:    300,
		Height:   169,
	}
	thumbnail := &discordgo.MessageEmbedThumbnail{
		URL:      services.BaseSiteURL + "t/p/w300" + tvShow.PosterPath,
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
			services.BaseSiteURL+"tv",
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

func GetEmbedForMovie(movie models.MovieResult) *discordgo.MessageEmbed {
	embedImage := &discordgo.MessageEmbedImage{
		URL:      services.BaseSiteURL + "t/p/w300" + movie.BackdropPath,
		ProxyURL: "",
		Width:    300,
		Height:   169,
	}
	thumbnail := &discordgo.MessageEmbedThumbnail{
		URL:      services.BaseSiteURL + "t/p/w300" + movie.PosterPath,
		ProxyURL: "",
		Width:    300,
		Height:   169,
	}

	releasedDate := movie.ReleaseDate
	if releasedDate == "" {
		releasedDate = "-"
	}

	embedFields := []*discordgo.MessageEmbedField{
		{
			Name:   "Status",
			Value:  movie.Status,
			Inline: true,
		},
		{
			Name:   "User Score",
			Value:  fmt.Sprintf("%.0f", movie.VoteAverage*10) + "%",
			Inline: true,
		},
		{
			Name:   "Released Date",
			Value:  releasedDate,
			Inline: true,
		},
	}

	if movie.Tagline != "" {
		embedFields = append(
			[]*discordgo.MessageEmbedField{
				{
					Name:   "Tagline",
					Value:  movie.Tagline,
					Inline: false,
				},
			},
			embedFields...,
		)
	}

	return MakeEmbed(
		fmt.Sprintf(
			"%v/%v",
			services.BaseSiteURL+"movie",
			movie.ID,
		),
		movie.Title,
		movie.Overview,
		embedImage,
		embedFields,
		thumbnail,
	)

}

func GetEmbedForCast(cast []models.Cast, contentId int, contentTitle, contentType string) *discordgo.MessageEmbed {

	var embedFields []*discordgo.MessageEmbedField

	for _, person := range cast[:20] {
		character := person.Character
		if person.Character == "" {
			character = "-"
		}
		embedFields = append(
			embedFields,
			&discordgo.MessageEmbedField{
				Name:   character,
				Value:  fmt.Sprintf("[%v](https://www.themoviedb.org/person/%v)", person.Name, person.ID),
				Inline: true,
			})
	}

	return MakeEmbed(
		fmt.Sprintf(
			"%v/%v",
			services.BaseSiteURL+contentType,
			contentId,
		),
		fmt.Sprintf("Cast for %v", contentTitle),
		"",
		nil,
		embedFields,
		nil,
	)

}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func MakeSearchResultTitles(results []interface{}) string {
	l := len(results)
	m := min(emojiLength, l)
	resultTitles := make([]string, l)
	var i string
	for index, result := range results[:m] {
		i = models.EmojiNumbersMap[index+1]
		resultTitles[index] = fmt.Sprintf("%s - %s", i, result)
	}
	for index, result := range results[m:] {
		i = strconv.Itoa(index + emojiLength + 1)
		resultTitles[index+emojiLength] = fmt.Sprintf("%s - %s", i, result)
	}
	return strings.Join(resultTitles, "\n")
}

func MakeTVShowSearchResultIdsMap(sr models.TVSearchResponse) map[int]int {
	idsMap := make(map[int]int)
	for index, result := range sr.Results {
		idsMap[index+1] = result.ID
	}
	return idsMap
}

func MakeMovieSearchResultIdsMap(sr models.MovieSearchResponse) map[int]int {
	idsMap := make(map[int]int)
	for index, result := range sr.Results {
		idsMap[index+1] = result.ID
	}
	return idsMap
}
