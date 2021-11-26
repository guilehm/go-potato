package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/guilehm/go-potato/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	tmdbAccessToken := os.Getenv("TMDB_ACCESS_TOKEN")
	if tmdbAccessToken == "" {
		fmt.Println("TMDB_ACCESS_TOKEN not set")
		return
	}

	discordToken := os.Getenv("DISCORD_TOKEN")
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error starting bot", err)
		return
	}

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	discord.AddHandler(handlers.MessageCreate)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = discord.Close()
	if err != nil {
		return
	}

}
