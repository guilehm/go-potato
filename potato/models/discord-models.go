package models

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type MessageData struct {
	MessageID  string `json:"message_id" bson:"message_id"`
	Text       string `json:"text" bson:"text"`
	Page       int    `json:"page" bson:"page"`
	TotalPages int    `json:"total_pages" bson:"total_pages"`
	Type       string `json:"type" bson:"type"`
}

type UserDiscord struct {
	discordgo.User `bson:",inline"`
	AvatarUrl      string    `json:"avatar_url" bson:"avatar_url"`
	DateChanged    time.Time `json:"date_changed" bson:"date_changed"`
	Likes          []int     `json:"likes" bson:"likes"`
}
