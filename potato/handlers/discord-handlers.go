package handlers

import (
	"fmt"
	"os"
	"strings"

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
			s.ChannelMessageSend(m.ChannelID, "Could not search: "+err.Error())
			return
		}

		if len(searchResponse.Results) == 0 {
			s.ChannelMessageSend(m.ChannelID, "Nothing found")
			return
		}

		titles := make([]string, len(searchResponse.Results))
		for index, result := range searchResponse.Results {
			titles[index] = result.Title
		}
		s.ChannelMessageSend(m.ChannelID, strings.Join(titles, "\n"))

	}

}
