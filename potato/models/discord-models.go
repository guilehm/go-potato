package models

type MessageData struct {
	MessageID string `json:"message_id" bson:"message_id"`
	Text      string `json:"text" bson:"text"`
	Page      int    `json:"page" bson:"page"`
}
