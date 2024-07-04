package model

import "time"

type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "active"
	ConversationStatusInactive ConversationStatus = "inactive"
)

type Conversation struct {
	ID        int                `json:"id"`
	AccountID int                `json:"account_id"`
	Account   *Account           `json:"account,omitempty"`
	Status    ConversationStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
