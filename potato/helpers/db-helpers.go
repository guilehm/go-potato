package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/guilehm/go-potato/potato/db"
	"github.com/guilehm/go-potato/potato/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UpdateTVShowDetail(tvShow models.TVShowResult, message *discordgo.Message) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}

	_, err := db.TVShowsCollection.UpdateOne(
		ctx, bson.M{"id": tvShow.ID}, bson.D{{Key: "$set", Value: tvShow}}, &opt,
	)
	if err != nil {
		fmt.Println("could not update TV Show #", tvShow.ID)
	}

	if err != nil {
		fmt.Println("Could not convert TV Show ID #", tvShow.ID)
		return
	}

	messageData := models.MessageData{
		MessageID:    message.ID,
		Type:         models.TD,
		ContentId:    tvShow.ID,
		ContentTitle: tvShow.Name,
	}
	_, err = db.MessagesDataCollection.InsertOne(ctx, messageData)
	if err != nil {
		fmt.Println("could save message data for tv-show #", tvShow.ID)
	}
}

func UpdateMovieDetail(movie models.MovieResult, message *discordgo.Message) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	upsert := true
	opt := options.UpdateOptions{Upsert: &upsert}

	_, err := db.MoviesCollection.UpdateOne(
		ctx, bson.M{"id": movie.ID}, bson.D{{Key: "$set", Value: movie}}, &opt,
	)
	if err != nil {
		fmt.Println("could not update Movie #", movie.ID)
	}

	if err != nil {
		fmt.Println("Could not convert Movie ID #", movie.ID)
		return
	}

	messageData := models.MessageData{
		MessageID:    message.ID,
		Type:         models.MD,
		ContentId:    movie.ID,
		ContentTitle: movie.Title,
	}
	_, err = db.MessagesDataCollection.InsertOne(ctx, messageData)
	if err != nil {
		fmt.Println("could save message data for movie #", movie.ID)
	}
}
