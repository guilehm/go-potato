package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connection() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("connected to MongoDB")
	return client
}

func DiscordUserCollection() *mongo.Collection {
	client := Connection()
	database := client.Database("go-potato")
	usersCollection := database.Collection("users")
	return usersCollection
}
