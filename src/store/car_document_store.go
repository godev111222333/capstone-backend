package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CarDocumentStore struct {
	db *gorm.DB
}

func NewCarDocumentStore(db *gorm.DB) *CarDocumentStore {
	return &CarDocumentStore{db: db}
}

func (s *CarDocumentStore) Create(carID int, document *model.Document) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(document).Error; err != nil {
			return err
		}
		return tx.Create(&model.CarDocument{
			CarID:      carID,
			DocumentID: document.ID,
		}).Error
	}); err != nil {
		fmt.Printf("CarDocumentStore: Create %v\n", err)
		return err
	}
	return nil
}
