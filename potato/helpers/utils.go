package helpers

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func MakeEmbed(url, title, description string) *discordgo.MessageEmbed {
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
	}
}
