package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/guilehm/go-potato/models"

	"github.com/bwmarrin/discordgo"
	"github.com/guilehm/go-potato/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func handleMyTVShowList(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.Find().SetLimit(100).SetProjection(
		bson.M{"_id": 0, "name": 1},
	)

	cur, err := db.TVShowsCollection.Find(
		ctx,
		bson.M{"id": bson.M{"$in": []int{888, 72705}}},
		opts,
	)

	if err != nil {
		fmt.Println("deu ruim ", err)
	}

	tvShows := make([]models.TVShowResult, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		var tvShow models.TVShowResult
		err := cur.Decode(&tvShow)
		if err != nil {
			fmt.Println("could not decode")
		}
		tvShows = append(tvShows, tvShow)
	}

}
