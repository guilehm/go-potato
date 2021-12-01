package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/guilehm/go-potato/models"

	"github.com/guilehm/go-potato/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/guilehm/go-potato/services"

	"github.com/bwmarrin/discordgo"
)

var service = services.TMDBService{AccessToken: os.Getenv("TMDB_ACCESS_TOKEN")}

func ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	switch r.Emoji.Name {
	case "⏭️":
		HandleNextPrev(s, r)
	case "⏮️":
		HandleNextPrev(s, r)
	case "❤️":
		HandleLikeAdd(s, r)
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

	if strings.ToLower(m.Content) == "hello" {
		handleHello(s, m)
		return
	}

	if strings.ToLower(m.Content) == "ping" {
		_ = s.ChannelTyping(m.ChannelID)
		_, err := s.ChannelMessageSend(m.ChannelID, "pong!")
		if err != nil {
			fmt.Println("could not send message for channel: " + m.ChannelID)
			fmt.Println(err.Error())
		}
		return
	}

	if strings.HasPrefix(m.Content, ".m ") {
		handleSearchMovies(s, m)
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

}

func handleHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user := models.UserDiscord{
		User:        *m.Author,
		AvatarUrl:   m.Author.AvatarURL(""),
		DateChanged: now,
		Likes:       []int{},
	}

	result, err := db.UsersCollection.UpdateOne(
		ctx, bson.M{"id": m.Author.ID}, bson.D{{Key: "$set", Value: &user}}, &opt,
	)

	var t string
	if result.UpsertedID != nil {
		t = "create"
	} else {
		t = "update"
	}
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Could not %v user", t))
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("User successfully %vd!", t))
	}

}
