package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleStockSearch(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)
	symbol := strings.Trim(m.Content[4:], " ")
	_, _ = s.ChannelMessageSend(m.ChannelID, "searching for "+symbol+"...")

	stock, err := stocksService.SearchStockPrice(symbol)
	if err != nil {
		_, _ = s.ChannelMessageSend(
			m.ChannelID,
			"Could not retrieve data for **"+symbol+"**. Please try again",
		)
		return
	}
	_ = s.ChannelTyping(m.ChannelID)
	_, _ = s.ChannelMessageSend(m.ChannelID, stock.CompanyName)

}
