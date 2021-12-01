package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func handlePing(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)
	_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
	if err != nil {
		fmt.Println("could not send message for channel: " + m.ChannelID)
		fmt.Println(err.Error())
	}
}
