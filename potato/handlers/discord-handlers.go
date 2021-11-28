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

	"github.com/guilehm/go-potato/helpers"
	"github.com/guilehm/go-potato/services"

	"github.com/bwmarrin/discordgo"
)

var service = services.TMDBService{AccessToken: os.Getenv("TMDB_ACCESS_TOKEN")}

func ReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}

	var m models.MessageData
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.MessagesDataCollection.FindOne(
		ctx, bson.M{"message_id": r.MessageID},
	).Decode(&m)
	if err != nil {
		fmt.Println("could not find message #" + r.MessageID)
		return
	}

	var page int
	var n int
	if r.Emoji.Name == "⏭️" {
		page = m.Page + 1
		n = 1
	} else if r.Emoji.Name == "⏮️" {
		page = m.Page - 1
		n = -1
	} else {
		return
	}

	_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", s.State.User.ID)
	_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", r.UserID)
	_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮", s.State.User.ID)
	_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮", r.UserID)

	var resultTitles string
	var title string
	if m.Type == models.T {
		searchResponse, err := service.SearchTvShows(m.Text, page)
		if err != nil {
			return
		}
		resultTitles = helpers.MakeTVShowSearchResultTitles(searchResponse)
		title = "TV Shows found:"

		if searchResponse.Page > 1 {
			_ = s.MessageReactionAdd(r.ChannelID, r.MessageID, "⏮️")
		}
		if searchResponse.Page < searchResponse.TotalPages {
			_ = s.MessageReactionAdd(r.ChannelID, r.MessageID, "⏭️")
		}
		if searchResponse.Page == searchResponse.TotalPages {
			_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", s.State.User.ID)
			_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", r.UserID)
		}
		if searchResponse.Page == 1 {
			_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮️", s.State.User.ID)
			_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮️", r.UserID)
		}
	} else if m.Type == models.M {
		searchResponse, err := service.SearchMovies(m.Text, page)
		if err != nil {
			return
		}
		title = "Movies found:"
		resultTitles = helpers.MakeMovieSearchResultTitles(searchResponse)

		if searchResponse.Page > 1 {
			_ = s.MessageReactionAdd(r.ChannelID, r.MessageID, "⏮️")
		}
		if searchResponse.Page < searchResponse.TotalPages {
			_ = s.MessageReactionAdd(r.ChannelID, r.MessageID, "⏭️")
		}
		if searchResponse.Page == searchResponse.TotalPages {
			_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", s.State.User.ID)
			_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", r.UserID)
		}
		if searchResponse.Page == 1 {
			_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮️", s.State.User.ID)
			_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮️", r.UserID)
		}
	} else {
		return
	}

	embed := helpers.MakeEmbed(
		"",
		title,
		resultTitles,
		&discordgo.MessageEmbedImage{},
		[]*discordgo.MessageEmbedField{},
		&discordgo.MessageEmbedThumbnail{},
	)
	_, err = s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed)
	if err != nil {
		return
	}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = db.MessagesDataCollection.UpdateOne(
			ctx,
			bson.M{"message_id": r.MessageID},
			bson.M{"$inc": bson.M{"page": n}},
		)
		if err != nil {
			return
		}
	}()

}

func ReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	_, _ = s.ChannelMessageEdit(r.ChannelID, r.MessageID, "overwritten")
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
	now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user := models.UserDiscord{
		User:        *m.Author,
		AvatarUrl:   m.Author.AvatarURL(""),
		DateChanged: now,
	}

	_, err := db.UsersCollection.UpdateOne(
		ctx, bson.M{"id": m.Author.ID}, bson.D{{Key: "$set", Value: &user}}, &opt,
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

	resultTitles := helpers.MakeMovieSearchResultTitles(searchResponse)
	message, err := s.ChannelMessageSendEmbed(
		m.ChannelID,
		helpers.MakeEmbed(
			"",
			"Movies found:",
			resultTitles,
			&discordgo.MessageEmbedImage{},
			[]*discordgo.MessageEmbedField{},
			&discordgo.MessageEmbedThumbnail{},
		),
	)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ops... Something weird happened")
	} else {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			upsert := true
			opt := options.UpdateOptions{Upsert: &upsert}

			messageData := models.MessageData{
				MessageID:  message.ID,
				Text:       text,
				Page:       1,
				TotalPages: searchResponse.TotalPages,
				Type:       models.M,
			}
			_, err = db.MessagesDataCollection.UpdateOne(
				ctx, bson.M{"message_id": message.ID}, bson.D{{Key: "$set", Value: &messageData}}, &opt,
			)
		}()
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

	resultTitles := helpers.MakeTVShowSearchResultTitles(searchResponse)
	message, err := s.ChannelMessageSendEmbed(
		m.ChannelID,
		helpers.MakeEmbed(
			"",
			"TV Shows found:",
			resultTitles,
			&discordgo.MessageEmbedImage{},
			[]*discordgo.MessageEmbedField{},
			&discordgo.MessageEmbedThumbnail{},
		),
	)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ops... Something weird happened")
	} else {
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			upsert := true
			opt := options.UpdateOptions{Upsert: &upsert}

			messageData := models.MessageData{
				MessageID:  message.ID,
				Text:       text,
				Page:       1,
				TotalPages: searchResponse.TotalPages,
				Type:       models.T,
			}
			_, err = db.MessagesDataCollection.UpdateOne(
				ctx, bson.M{"message_id": message.ID}, bson.D{{Key: "$set", Value: &messageData}}, &opt,
			)
		}()
	}

	if searchResponse.Page < searchResponse.TotalPages {
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "⏭️")
	}
	if searchResponse.Page > 1 {
		_ = s.MessageReactionAdd(m.ChannelID, message.ID, "⏮️")
	}

}
