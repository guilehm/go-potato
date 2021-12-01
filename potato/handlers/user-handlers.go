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

func handleMyTVShowList(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.UserDiscord
	if err := db.UsersCollection.FindOne(ctx, bson.M{"id": m.Author.ID}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			_, _ = s.ChannelMessageSend(m.ChannelID, "Please say \"hello\" to create your user")
		}
		return
	}

	opts := options.Find()
	cur, err := db.TVShowsCollection.Find(
		ctx,
		bson.M{"id": bson.M{"$in": user.Likes}},
		opts,
	)

	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Could not find your list")
	}

	tvShows := make([]models.TVShowResult, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		var tvShow models.TVShowResult
		err := cur.Decode(&tvShow)
		if err != nil {
			fmt.Println("could not decode")
		}
		tvShows = append(tvShows, tvShow)
		go func() {
			_, _ = s.ChannelMessageSendEmbed(m.ChannelID, helpers.GetSimpleEmbedForTVShow(tvShow))
		}()
	}

}
