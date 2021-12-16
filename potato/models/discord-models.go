package models

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	WHITE              int = 16777215
	BLURPLE            int = 5793266
	GREYPLE            int = 10070709
	DARK_BUT_NOT_BLACK int = 2895667
	NOT_QUITE_BLACK    int = 2303786
	GREEN              int = 5763719
	YELLOW             int = 16705372
	FUSCHIA            int = 15418782
	RED                int = 15548997
	BLACK              int = 2303786
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
