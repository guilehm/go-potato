package models

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	Blue      int = 3447003
	DarkBlue  int = 2123412
	Green     int = 3066993
	DarkGreen int = 2067276
	Fuschia   int = 15418782
	Red       int = 15548997
	Black     int = 2303786
	Yellow    int = 16776960
)

var EmojiNumbersMap = map[int]string{
	1: "1️⃣",
	2: "2️⃣",
	3: "3️⃣",
	4: "4️⃣",
	5: "5️⃣",
}

type MessageData struct {
	MessageID    string      `json:"message_id" bson:"message_id"`
	Text         string      `json:"text" bson:"text"`
	Page         int         `json:"page" bson:"page"`
	TotalPages   int         `json:"total_pages" bson:"total_pages"`
	Type         string      `json:"type" bson:"type"`
	ContentId    int         `json:"content_id" bson:"content_id"`
	ContentTitle string      `json:"content_title" bson:"content_title"`
	IDsMap       map[int]int `json:"ids_map" bson:"ids_map"`
}

type UserDiscord struct {
	discordgo.User `bson:",inline"`
	AvatarUrl      string    `json:"avatar_url" bson:"avatar_url"`
	DateChanged    time.Time `json:"date_changed" bson:"date_changed"`
	Likes          []int     `json:"likes" bson:"likes"`
}
