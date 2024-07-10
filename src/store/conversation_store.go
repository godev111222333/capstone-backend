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

func (s *ConversationStore) Get(offset, limit int) ([]*model.Conversation, error) {
	var res []*model.Conversation
	if limit == 0 {
		limit = 10000
	}

	if err := s.db.Order("id desc").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
		fmt.Printf("ConversationStore: Get %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *ConversationStore) GetByAccID(acctID int) (*model.Conversation, error) {
	res := &model.Conversation{}
	row := s.db.
		Where("account_id = ? and status = ?", acctID, string(model.ConversationStatusActive)).
		Find(res)
	if err := row.Error; err != nil {
		fmt.Printf("ConversationStore: GetByAccID %v\n", err)
		return nil, err
	}

	if row.RowsAffected == 0 {
		return nil, nil
	}

	return res, nil
}

func (s *ConversationStore) GetByID(ID int) (*model.Conversation, error) {
	res := &model.Conversation{}
	row := s.db.
		Where("id = ? and status = ?", ID, string(model.ConversationStatusActive)).
		Preload("Account").
		Find(res)
	if err := row.Error; err != nil {
		fmt.Printf("ConversationStore: GetByID %v\n", err)
		return nil, err
	}

	if row.RowsAffected == 0 {
		return nil, nil
	}

	return res, nil
}
