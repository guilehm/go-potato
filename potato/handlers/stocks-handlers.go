package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func handleStockSearch(s *discordgo.Session, m *discordgo.MessageCreate) {
	symbol := strings.Trim(m.Content[4:], " ")
	_, _ = s.ChannelMessageSend(m.ChannelID, "searching for "+symbol)
}
