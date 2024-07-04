package model

import "time"

type Message struct {
	ID             int           `json:"id"`
	ConversationID int           `json:"conversation_id"`
	Conversation   *Conversation `json:"conversation"`
	Sender         int           `json:"sender"`
	Content        string        `json:"content"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
