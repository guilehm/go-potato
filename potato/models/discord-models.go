package models

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type MessageData struct {
	MessageID    string `json:"message_id" bson:"message_id"`
	Text         string `json:"text" bson:"text"`
	Page         int    `json:"page" bson:"page"`
	TotalPages   int    `json:"total_pages" bson:"total_pages"`
	Type         string `json:"type" bson:"type"`
	ContentId    int    `json:"content_id" bson:"content_id"`
	ContentTitle string `json:"content_title" bson:"content_title"`
}

type UserDiscord struct {
	discordgo.User `bson:",inline"`
	AvatarUrl      string    `json:"avatar_url" bson:"avatar_url"`
	DateChanged    time.Time `json:"date_changed" bson:"date_changed"`
	Likes          []int     `json:"likes" bson:"likes"`
}
