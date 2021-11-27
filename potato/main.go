package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/guilehm/go-potato/handlers"

	"github.com/bwmarrin/discordgo"
)

func main() {
	tmdbAccessToken := os.Getenv("TMDB_ACCESS_TOKEN")
	if tmdbAccessToken == "" {
		log.Fatal("TMDB_ACCESS_TOKEN not set")
	}

	discordToken := os.Getenv("DISCORD_TOKEN")
	discord, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal("error starting bot", err)
	}

	err = discord.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}
	discord.AddHandler(handlers.MessageCreate)
	discord.AddHandler(handlers.ReactionAdd)
	discord.AddHandler(handlers.ReactionRemove)

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = discord.Close()
	if err != nil {
		return
	}

}
