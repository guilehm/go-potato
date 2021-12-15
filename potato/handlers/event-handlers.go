package handlers

import (
	"os"
	"strings"

	"github.com/guilehm/go-potato/potato/services"

	"github.com/bwmarrin/discordgo"
)

var service = services.TMDBService{AccessToken: os.Getenv("TMDB_ACCESS_TOKEN")}
var stocksService = services.StocksService{SecretKey: os.Getenv("STOCKS_API_SECRET_KEY")}

func ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	switch r.Emoji.Name {
	case "⏮️", "⏭️":
		HandleNextPrev(s, r)
	case "❤️":
		HandleLikeAdd(s, r)
	case "1️⃣", "2️⃣", "3️⃣":
		HandleNumberAdd(s, r)
	case "👪":
		HandleCastingAdd(s, r)
	}
}

func ReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == s.State.User.ID {
		return
	}

	switch r.Emoji.Name {
	case "❤️":
		HandleLikeRemove(s, r)
	}
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.ToLower(m.Content) == "ping" {
		handlePing(s, m)
		return
	}

	if strings.ToLower(m.Content) == "hello" {
		handleHello(s, m)
		return
	}

	if strings.HasPrefix(m.Content, ".m ") {
		handleSearchMovies(s, m)
		return
	}

	if strings.HasPrefix(m.Content, ".md ") {
		handleMovieDetail(s, m)
		return
	}

	if strings.HasPrefix(m.Content, ".t ") {
		handleSearchTVShows(s, m)
		return
	}

	if strings.HasPrefix(m.Content, ".td ") {
		handleTVShowDetail(s, m)
		return
	}

	if strings.HasPrefix(m.Content, ".tl") {
		handleTVShowLikeList(s, m)
		return
	}

	if strings.HasPrefix(m.Content, ".st") {
		handleStockSearch(s, m)
		return
	}

}
