package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/guilehm/go-potato/db"
	"github.com/guilehm/go-potato/helpers"
	"github.com/guilehm/go-potato/models"
	"go.mongodb.org/mongo-driver/bson"
)

func HandleNextPrev(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
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
	}

	_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", s.State.User.ID)
	_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", r.UserID)
	_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮", s.State.User.ID)
	_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮", r.UserID)

	var resultTitles string
	var title string
	var srPage int
	var srTotalPages int
	if m.Type == models.T {
		searchResponse, err := service.SearchTvShows(m.Text, page)
		if err != nil {
			return
		}
		resultTitles = helpers.MakeTVShowSearchResultTitles(searchResponse)
		title = "TV Shows found:"
		srPage = searchResponse.Page
		srTotalPages = searchResponse.TotalPages
	} else if m.Type == models.M {
		searchResponse, err := service.SearchMovies(m.Text, page)
		if err != nil {
			return
		}
		title = "Movies found:"
		resultTitles = helpers.MakeMovieSearchResultTitles(searchResponse)
		srPage = searchResponse.Page
		srTotalPages = searchResponse.TotalPages
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

	if srPage > 1 {
		_ = s.MessageReactionAdd(r.ChannelID, r.MessageID, "⏮️")
	}
	if srPage < srTotalPages {
		_ = s.MessageReactionAdd(r.ChannelID, r.MessageID, "⏭️")
	}
	if srPage == srTotalPages {
		_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", s.State.User.ID)
		_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏭️", r.UserID)
	}
	if srPage == 1 {
		_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮️", s.State.User.ID)
		_ = s.MessageReactionRemove(r.ChannelID, r.MessageID, "⏮️", r.UserID)
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
