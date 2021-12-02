package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

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

	_ = s.MessageReactionsRemoveAll(r.ChannelID, r.MessageID)

	var resultTitles string
	var title string
	var rCount int
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
		rCount = len(searchResponse.Results)
	} else if m.Type == models.M {
		searchResponse, err := service.SearchMovies(m.Text, page)
		if err != nil {
			return
		}
		title = "Movies found:"
		resultTitles = helpers.MakeMovieSearchResultTitles(searchResponse)
		srPage = searchResponse.Page
		srTotalPages = searchResponse.TotalPages
		rCount = len(searchResponse.Results)
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
	msg, err := s.ChannelMessageEditEmbed(r.ChannelID, r.MessageID, embed)
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

	for i := 1; i <= rCount && i <= 3; i++ {
		_ = s.MessageReactionAdd(msg.ChannelID, msg.ID, models.EmojiNumbersMap[i])
	}

}

func HandleLikeAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.UserDiscord
	if err := db.UsersCollection.FindOne(ctx, bson.M{"id": r.UserID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, _ = s.ChannelMessageSend(r.ChannelID, "Please say \"hello\" to create your user")
		}
		return
	}

	var message models.MessageData
	if err := db.MessagesDataCollection.FindOne(
		ctx,
		bson.M{"message_id": r.MessageID},
	).Decode(&message); err != nil {
		fmt.Printf("Could not find message #%v\n", r.MessageID)
		return
	}

	_, err := db.UsersCollection.UpdateOne(
		ctx,
		bson.M{"id": user.ID},
		bson.M{"$addToSet": bson.M{"likes": message.ContentId}},
	)
	if err != nil {
		fmt.Printf("Could not add like to user #%v for message #%v. %v\n", r.UserID, message.MessageID, err)
		return
	}

	_, _ = s.ChannelMessageSend(
		r.ChannelID, fmt.Sprintf(
			"<@%v> \"%v\" successfully **added** to your like list",
			user.ID,
			message.ContentTitle,
		),
	)
}

func HandleLikeRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.UserDiscord
	if err := db.UsersCollection.FindOne(ctx, bson.M{"id": r.UserID}).Decode(&user); err != nil {
		return
	}

	var message models.MessageData
	if err := db.MessagesDataCollection.FindOne(
		ctx,
		bson.M{"message_id": r.MessageID},
	).Decode(&message); err != nil {
		fmt.Printf("Could not find message #%v\n", r.MessageID)
		return
	}

	_, err := db.UsersCollection.UpdateOne(
		ctx,
		bson.M{"id": user.ID},
		bson.M{"$pull": bson.M{"likes": message.ContentId}},
	)
	if err != nil {
		fmt.Printf("Could not remove like to user #%v for message #%v\n", r.UserID, message.MessageID)
		return
	}

	_, _ = s.ChannelMessageSend(
		r.ChannelID, fmt.Sprintf(
			"<@%v> \"%v\" successfully **removed** from your like list",
			user.ID,
			message.ContentTitle,
		),
	)
}

func HandleNumberAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message := models.MessageData{}
	if err := db.MessagesDataCollection.FindOne(
		ctx,
		bson.M{"message_id": r.MessageID},
	).Decode(&message); err != nil {
		fmt.Printf("Could not find message #%v\n", r.MessageID)
		return
	}

	emojiIndex := 0
	for key, value := range models.EmojiNumbersMap {
		if r.Emoji.Name == value {
			emojiIndex = key
			break
		}
	}

	// TODO: Add condition for movies
	if message.Type == models.T {
		intTVShowID := message.IDsMap[emojiIndex]
		tvShowID := strconv.Itoa(intTVShowID)
		tvShow, err := service.GetTVShowDetail(tvShowID)
		if err != nil {
			_, _ = s.ChannelMessageSend(r.ChannelID, "Could not get tv show detail: "+err.Error())
			return
		}

		msg, err := s.ChannelMessageSendEmbed(
			r.ChannelID,
			helpers.GetEmbedForTVShow(tvShow),
		)

		go func() {
			helpers.UpdateTVShowDetail(tvShow, msg)
		}()

		if err != nil {
			_, _ = s.ChannelMessageSend(r.ChannelID, "Ops... Something weird happened")
		}
		_ = s.MessageReactionAdd(r.ChannelID, msg.ID, "❤️")

		messageData := models.MessageData{
			MessageID:    msg.ID,
			Type:         models.TD,
			ContentId:    intTVShowID,
			ContentTitle: tvShow.Name,
		}
		_, err = db.MessagesDataCollection.InsertOne(ctx, messageData)
		if err != nil {
			fmt.Println("could save message data for tv-show #" + tvShowID)
		}
	}

}
