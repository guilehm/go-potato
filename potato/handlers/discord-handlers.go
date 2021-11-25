package handlers

import (
	"fmt"
	"os"
	"strings"
	"time"

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
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
			URL:         "",
			Type:        "",
			Title:       "Movies found:",
			Description: strings.Join(titles, "\n"),
			Timestamp:   time.Now().Format("2006-01-02 15:04"),
			Color:       3447003,
			Footer: &discordgo.MessageEmbedFooter{
				IconURL:      "",
				Text:         "go potato",
				ProxyIconURL: "",
			},
			Image:     nil,
			Thumbnail: nil,
			Video:     nil,
			Provider:  nil,
			Author: &discordgo.MessageEmbedAuthor{
				Name:         "the movie db",
				IconURL:      "https://www.themoviedb.org/assets/2/apple-touch-icon-57ed4b3b0450fd5e9a0c20f34e814b82adaa1085c79bdde2f00ca8787b63d2c4.png",
				URL:          "https://www.themoviedb.org/",
				ProxyIconURL: "",
			},
			Fields: nil,
		})

	}

}
