package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/guilehm/go-potato/db"
	"github.com/guilehm/go-potato/helpers"
	"github.com/guilehm/go-potato/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

func handleTVShowDetail(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)

	tvShowID := strings.Trim(m.Content[4:], " ")
	tvShow, err := service.GetTVShowDetail(tvShowID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Could not get tv show detail: "+err.Error())
		return
	}

	message, err := s.ChannelMessageSendEmbed(
		m.ChannelID,
		helpers.GetEmbedForTVShow(tvShow),
	)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ops... Something weird happened")
	}
	_ = s.MessageReactionAdd(m.ChannelID, message.ID, "❤️")

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

		messageData := models.MessageData{
			MessageID:    message.ID,
			Type:         models.TD,
			ContentId:    tvShowID,
			ContentTitle: tvShow.Name,
		}
		_, err = db.MessagesDataCollection.InsertOne(ctx, messageData)
		if err != nil {
			fmt.Println("could save message data for tv-show #" + tvShowID)
		}
	}()
}
