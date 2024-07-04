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
func (s *MessageStore) GetByConversationID(convID, offset, limit int) ([]*model.Message, error) {
	var res []*model.Message
	if limit == 0 {
		limit = 10000
	}

	if err := s.db.Where("conversation_id = ?", convID).Order("id desc").Offset(offset).Limit(limit).Preload("Account").Find(&res).Error; err != nil {
		fmt.Printf("MessageStore: GetByConversationID %v\n", err)
		return nil, err
	}

	return res, nil
}
