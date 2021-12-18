package handlers

import (
	"errors"
	"strings"

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
	_, _ = s.ChannelMessageSend(m.ChannelID, stock.CompanyName)

}
