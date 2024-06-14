package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type DocumentStore struct {
	db *gorm.DB
}

func NewDocumentStore(db *gorm.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

func (s *DocumentStore) Create(m *model.Document) error {
	if err := s.db.Create(m).Error; err != nil {
		fmt.Printf("DocumentStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *DocumentStore) GetByCategory(
	accID int,
	category model.DocumentCategory,
	limit int,
) ([]*model.Document, error) {
	var docs []*model.Document
	if err := s.db.Where(
		"account_id = ? and category = ? and status = ?",
		accID,
		string(category),
		model.DocumentStatusActive,
	).Order("id desc").Limit(limit).Find(&docs).Error; err != nil {
		fmt.Printf("DocumentStore: GetByCategory %v\n", err)
		return nil, err
	}

	reverse := make([]*model.Document, len(docs))
	for i := range docs {
		reverse[i] = docs[len(docs)-i-1]
	}

	return reverse, nil
}
