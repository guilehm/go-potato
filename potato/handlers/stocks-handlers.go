package handlers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/guilehm/go-potato/potato/services"

	"github.com/bwmarrin/discordgo"
)

func handleStockSearch(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)
	symbol := strings.Trim(m.Content[4:], " ")
	_, _ = s.ChannelMessageSend(m.ChannelID, "searching for "+symbol+"...")

	stock, err := stocksService.SearchStockPrice(symbol)
	if err != nil {
		errorMessage := "Could not retrieve data for **" + symbol + "**. Please try again"
		if errors.Is(err, services.ErrApiNotSet) {
			errorMessage = "Ops... Api Key not set!"
		}
		_, _ = s.ChannelMessageSend(
			m.ChannelID,
			errorMessage,
		)
		return
	}
	_ = s.ChannelTyping(m.ChannelID)
	_, _ = s.ChannelMessageSendEmbed(
		m.ChannelID,
		&discordgo.MessageEmbed{
			URL:         stock.Website,
			Title:       stock.CompanyName,
			Description: stock.Description,
			Timestamp:   time.Now().Format("2006-01-02 15:04"),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Symbol",
					Value:  stock.Symbol,
					Inline: true,
				},
				{
					Name:   "Document",
					Value:  stock.Document,
					Inline: true,
				},
				{
					Name:   "Price",
					Value:  fmt.Sprintf("%.2f %s", stock.Price, stock.Currency),
					Inline: true,
				},
			},
		},
	)

}
