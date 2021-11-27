package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/guilehm/go-potato/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/guilehm/go-potato/helpers"
	"github.com/guilehm/go-potato/services"

	"github.com/bwmarrin/discordgo"
)

var service = services.TMDBService{AccessToken: os.Getenv("TMDB_ACCESS_TOKEN")}

func ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
}

func ReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
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

func handleTVShowDetail(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)

	tvShowID := strings.Trim(m.Content[4:], " ")
	tvShow, err := service.GetTVShowDetail(tvShowID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Could not get tv show detail: "+err.Error())
		return
	}

	_, err = s.ChannelMessageSendEmbed(
		m.ChannelID,
		helpers.GetEmbedForTVShow(tvShow),
	)

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ops... Something weird happened")
	}

	go func() {

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}

		_, err := db.TVShowsCollection.UpdateOne(
			ctx, bson.M{"id": tvShow.ID}, bson.D{{Key: "$set", Value: tvShow}}, &opt,
		)
		if err != nil {
			fmt.Println("could not update TV Show #" + tvShowID)
		}
	}()
}

func handleHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}

	// TODO: add avatar_url field (m.Author.AvatarURL(""))
	_, err := db.UsersCollection.UpdateOne(
		ctx, bson.M{"id": m.Author.ID}, bson.D{{Key: "$set", Value: m.Author}}, &opt,
	)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Could not update user: "+err.Error())
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, "User successfully updated!")
	}

}

func handleSearchMovies(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)

	text := strings.Trim(m.Content[3:], " ")
	searchResponse, err := service.SearchMovies(text, 1)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Could not search movie: "+err.Error())
		return
	}

	if len(searchResponse.Results) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, `Nothing found for "`+text+`"`)
		return
	}

	resultTitles := make([]string, len(searchResponse.Results))
	for index, result := range searchResponse.Results {
		resultTitles[index] = fmt.Sprintf("%v *(%v)*", result.Title, result.ID)
	}
	message, err := s.ChannelMessageSendEmbed(
		m.ChannelID,
		helpers.MakeEmbed(
			"",
			"Movies found:",
			strings.Join(resultTitles, "\n"),
			&discordgo.MessageEmbedImage{},
			[]*discordgo.MessageEmbedField{},
			&discordgo.MessageEmbedThumbnail{},
		),
	)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ops... Something weird happened")
	}

	if searchResponse.Page < searchResponse.TotalPages {
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "⏭️")
	}
	if searchResponse.Page > 1 {
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "⏮️")
	}
}

func handleSearchTVShows(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)

	text := strings.Trim(m.Content[3:], " ")
	searchResponse, err := service.SearchTvShows(text, 1)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Could not search Tv Shows: "+err.Error())
		return
	}

	if len(searchResponse.Results) == 0 {
		_, _ = s.ChannelMessageSend(m.ChannelID, `Nothing found for "`+text+`"`)
		return
	}

	resultTitles := make([]string, len(searchResponse.Results))
	for index, result := range searchResponse.Results {
		resultTitles[index] = fmt.Sprintf("%v *(%v)*", result.Name, result.ID)
	}
	message, err := s.ChannelMessageSendEmbed(
		m.ChannelID,
		helpers.MakeEmbed(
			"",
			"TV Shows found:",
			strings.Join(resultTitles, "\n"),
			&discordgo.MessageEmbedImage{},
			[]*discordgo.MessageEmbedField{},
			&discordgo.MessageEmbedThumbnail{},
		),
	)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ops... Something weird happened")
	}

	if searchResponse.Page < searchResponse.TotalPages {
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "⏭️")
	}
	if searchResponse.Page > 1 {
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "⏮️")
	}

}
