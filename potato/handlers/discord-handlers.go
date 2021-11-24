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
	_ = s.ChannelTyping(m.ChannelID)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
		if err != nil {
			fmt.Println("could not send message for channel: " + m.ChannelID)
			fmt.Println(err.Error())
		}
	}

	if strings.HasPrefix(m.Content, ".s") {
		text := strings.Trim(m.Content[3:], " ")
		s.ChannelMessageSend(m.ChannelID, "searching for: "+text)
		searchResponse, err := service.SearchMovie(text)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Could not search. Please try again later."+err.Error())
		}
		for _, result := range searchResponse.Results {
			s.ChannelMessage(m.ChannelID, result.Title)
		}
	}

}
