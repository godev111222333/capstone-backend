package model

import "time"

type NotificationStatus string

const (
	NotificationStatusActive   NotificationStatus = "active"
	NotificationStatusInactive NotificationStatus = "inactive"
)

type Notification struct {
	ID        int                `json:"id"`
	AccountID int                `json:"account_id"`
	Account   *Account           `json:"account"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	URL       string             `json:"url"`
	Status    NotificationStatus `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
