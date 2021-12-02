package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/guilehm/go-potato/helpers"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/guilehm/go-potato/models"

	"github.com/bwmarrin/discordgo"
	"github.com/guilehm/go-potato/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func handleHello(s *discordgo.Session, m *discordgo.MessageCreate) {
	_ = s.ChannelTyping(m.ChannelID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	var user models.UserDiscord
	err := db.UsersCollection.FindOne(ctx, bson.M{"id": m.Author.ID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			user.User = *m.Author
			user.AvatarUrl = m.Author.AvatarURL("")
			user.DateChanged = now
			user.Likes = []int{}

			_, err := db.UsersCollection.InsertOne(ctx, user)
			if err != nil {
				_, _ = s.ChannelMessageSend(
					m.ChannelID,
					fmt.Sprintf("Could not create your user. Please try again."),
				)
				return
			}
			_, _ = s.ChannelMessageSend(m.ChannelID, "User successfully created!")
			return
		}
		_, _ = s.ChannelMessageSend(
			m.ChannelID,
			"An error occurred while trying to create your user. Please try again.",
		)
		return
	}

	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}
	user = models.UserDiscord{
		User:        *m.Author,
		AvatarUrl:   m.Author.AvatarURL(""),
		DateChanged: now,
		Likes:       user.Likes,
	}

	_, err = db.UsersCollection.UpdateOne(
		ctx, bson.M{"id": m.Author.ID}, bson.D{{Key: "$set", Value: &user}}, &opt,
	)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Could not update your user")
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, "User successfully updated!")
	}

}

func handleTVShowLikeList(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.UserDiscord
	if err := db.UsersCollection.FindOne(ctx, bson.M{"id": m.Author.ID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Please say \"hello\" to create your user")
		}
		return
	}

	_, _ = s.ChannelMessageSend(
		m.ChannelID, fmt.Sprintf(
			"<@%v> Here is your TV Show like list.",
			user.ID,
		),
	)

	opts := options.Find()
	cur, err := db.TVShowsCollection.Find(
		ctx,
		bson.M{"id": bson.M{"$in": user.Likes}},
		opts,
	)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Could not find your list")
		return
	}

	tvShows := make([]models.TVShowResult, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		var tvShow models.TVShowResult
		err := cur.Decode(&tvShow)
		if err != nil {
			fmt.Println("could not decode tv-show")
			continue
		}
		tvShows = append(tvShows, tvShow)
		go func() {
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, helpers.GetSimpleEmbedForTVShow(tvShow))
		}()
	}

}
