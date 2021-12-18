package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bwmarrin/discordgo"
	"github.com/guilehm/go-potato/potato/db"
	"github.com/guilehm/go-potato/potato/helpers"
	"github.com/guilehm/go-potato/potato/models"
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

	var (
		resultTitles string
		title        string
		rCount       int
		srPage       int
		srTotalPages int
		idsMap       map[int]int
		color        int
	)

	if m.Type == models.T {
		searchResponse, err := service.SearchTvShows(m.Text, page)
		if err != nil {
			return
		}

		results := make([]interface{}, len(searchResponse.Results))
		for i, result := range searchResponse.Results {
			results[i] = result
		}
		resultTitles = helpers.MakeSearchResultTitles(results)

		title = "TV Shows found:"
		srPage = searchResponse.Page
		srTotalPages = searchResponse.TotalPages
		rCount = len(searchResponse.Results)
		idsMap = helpers.MakeTVShowSearchResultIdsMap(searchResponse)
		color = models.Blue
	} else if m.Type == models.M {
		searchResponse, err := service.SearchMovies(m.Text, page)
		if err != nil {
			return
		}
		title = "Movies found:"

		results := make([]interface{}, len(searchResponse.Results))
		for i, result := range searchResponse.Results {
			results[i] = result
		}
		resultTitles = helpers.MakeSearchResultTitles(results)

		srPage = searchResponse.Page
		srTotalPages = searchResponse.TotalPages
		rCount = len(searchResponse.Results)
		idsMap = helpers.MakeMovieSearchResultIdsMap(searchResponse)
		color = models.Green
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
		color,
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
			bson.M{
				"$inc": bson.M{"page": n},
				"$set": bson.M{"ids_map": idsMap},
			},
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
		if err != nil {
			_, _ = s.ChannelMessageSend(r.ChannelID, "Ops... Something weird happened")
		}
		_ = s.MessageReactionAdd(r.ChannelID, msg.ID, "❤️")
		_ = s.MessageReactionAdd(r.ChannelID, msg.ID, "👪")

		go func() {
			helpers.UpdateTVShowDetail(tvShow, msg)
		}()

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

	if message.Type == models.M {
		intMovieID := message.IDsMap[emojiIndex]
		movieID := strconv.Itoa(intMovieID)
		movie, err := service.GetMovieDetail(movieID)
		if err != nil {
			_, _ = s.ChannelMessageSend(r.ChannelID, "Could not get movie detail: "+err.Error())
			return
		}

		msg, err := s.ChannelMessageSendEmbed(
			r.ChannelID,
			helpers.GetEmbedForMovie(movie),
		)

		if err != nil {
			_, _ = s.ChannelMessageSend(r.ChannelID, "Ops... Something weird happened")
			fmt.Println(err)
		}
		_ = s.MessageReactionAdd(r.ChannelID, msg.ID, "❤️")
		_ = s.MessageReactionAdd(r.ChannelID, msg.ID, "👪")

		go func() {
			helpers.UpdateMovieDetail(movie, msg)
		}()

		messageData := models.MessageData{
			MessageID:    msg.ID,
			Type:         models.MD,
			ContentId:    intMovieID,
			ContentTitle: movie.Title,
		}
		_, err = db.MessagesDataCollection.InsertOne(ctx, messageData)
		if err != nil {
			fmt.Println("could save message data for movie #" + movieID)
		}
	}

}

func HandleCastingAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

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

	if message.Type == models.MD {
		var movie models.MovieResult
		if err := db.MoviesCollection.FindOne(
			ctx,
			bson.M{"id": message.ContentId},
		).Decode(&movie); err != nil {
			_, _ = s.ChannelMessageSend(r.ChannelID, "Could not find the requested movie")
			return
		}

		_, err := s.ChannelMessageSendEmbed(
			r.ChannelID,
			helpers.GetEmbedForCast(movie.Credits.Cast, movie.ID, movie.Title, "movie"),
		)
		if err != nil {
			fmt.Println("could not send message for channel: " + r.ChannelID)
			fmt.Println(err.Error())
		}

	}

	if message.Type == models.TD {

		var tvShow models.TVShowResult
		if err := db.TVShowsCollection.FindOne(
			ctx,
			bson.M{"id": message.ContentId},
		).Decode(&tvShow); err != nil {
			_, _ = s.ChannelMessageSend(r.ChannelID, "Could not find the requested tv show")
			return
		}

		_, err := s.ChannelMessageSendEmbed(
			r.ChannelID,
			helpers.GetEmbedForCast(tvShow.Credits.Cast, tvShow.ID, tvShow.Name, "tv"),
		)
		if err != nil {
			fmt.Println("could not send message for channel: " + r.ChannelID)
			fmt.Println(err.Error())
		}

	}

}
