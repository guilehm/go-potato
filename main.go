package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/guilehm/go-potato/potato/handlers"

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

	fmt.Println("Bot is now running.")

	err = discord.Close()
	if err != nil {
		return
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Hi there!")
	})
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))

}
