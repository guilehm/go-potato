package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
		if err != nil {
			fmt.Println("could not send message for channel: " + m.ChannelID)
		}
	}

}

func main() {

	discordToken := os.Getenv("DISCORD_TOKEN")
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error starting bot", err)
		return
	}

	discord.AddHandler(messageCreate)

}
