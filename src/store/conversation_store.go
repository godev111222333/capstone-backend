package store

import (
	"fmt"

	"github.com/godev111222333/capstone-backend/src/model"
	"gorm.io/gorm"
)

type ConversationStore struct {
	db *gorm.DB
}

func NewConversationStore(db *gorm.DB) *ConversationStore {
	return &ConversationStore{db: db}
}

func (s *ConversationStore) Create(c *model.Conversation) error {
	if err := s.db.Create(c).Error; err != nil {
		fmt.Printf("ConversationStore: Create %v\n", err)
		return err
	}

	return nil
}
