package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/guilehm/go-potato/helpers"
	"github.com/guilehm/go-potato/services"

	"github.com/bwmarrin/discordgo"
)

var service = services.TMDBService{AccessToken: os.Getenv("TMDB_ACCESS_TOKEN")}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_ = s.ChannelTyping(m.ChannelID)
		_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
		if err != nil {
			fmt.Println("could not send message for channel: " + m.ChannelID)
			fmt.Println(err.Error())
		}
	}

	if strings.HasPrefix(m.Content, ".s") {
		_ = s.ChannelTyping(m.ChannelID)

		text := strings.Trim(m.Content[3:], " ")
		searchResponse, err := service.SearchMovie(text)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not search movies: "+err.Error())
			return
		}

		if len(searchResponse.Results) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Nothing found for "+text)
			return
		}

		resultTitles := make([]string, len(searchResponse.Results))
		for index, result := range searchResponse.Results {
			resultTitles[index] = result.Title
		}
		_, err = s.ChannelMessageSendEmbed(
			m.ChannelID,
			helpers.MakeEmbed(
				"",
				"Movies found:",
				strings.Join(resultTitles, "\n"),
			),
		)
	}

	if strings.HasPrefix(m.Content, ".s") {
		_ = s.ChannelTyping(m.ChannelID)

		text := strings.Trim(m.Content[3:], " ")
		searchResponse, err := service.SearchTvShows(text)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not search Tv Shows: "+err.Error())
			return
		}

		if len(searchResponse.Results) == 0 {
			s.ChannelMessageSend(m.ChannelID, `Nothing found for "`+text+`"`)
			return
		}

		resultTitles := make([]string, len(searchResponse.Results))
		for index, result := range searchResponse.Results {
			resultTitles[index] = result.Name
		}
		_, err = s.ChannelMessageSendEmbed(
			m.ChannelID,
			helpers.MakeEmbed(
				"",
				"TV Shows found:",
				strings.Join(resultTitles, "\n"),
			),
		)
	}

}
