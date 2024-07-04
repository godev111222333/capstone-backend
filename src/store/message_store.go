package store

import (
	"fmt"

	"github.com/godev111222333/capstone-backend/src/model"
	"gorm.io/gorm"
)

type MessageStore struct {
	db *gorm.DB
}

func NewMessageStore(db *gorm.DB) *MessageStore {
	return &MessageStore{db: db}
}

func (s *MessageStore) Create(m *model.Message) error {
	if err := s.db.Create(m).Error; err != nil {
		fmt.Printf("MessageStore: Create %v\n", err)
		return err
	}
	return nil
}
