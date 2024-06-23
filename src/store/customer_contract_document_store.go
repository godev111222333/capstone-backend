package store

import (
	"fmt"
	"github.com/godev111222333/capstone-backend/src/model"
	"gorm.io/gorm"
)

type CustomerContractDocumentStore struct {
	db *gorm.DB
}

func NewCustomerContractDocumentStore(db *gorm.DB) *CustomerContractDocumentStore {
	return &CustomerContractDocumentStore{db: db}
}

func (s *CustomerContractDocumentStore) Create(cusContractID int, docs []*model.Document) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(docs).Error; err != nil {
			return err
		}

		for _, doc := range docs {
			if err := tx.Create(&model.CustomerContractDocument{
				CustomerContractID: cusContractID,
				DocumentID:         doc.ID,
			}).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fmt.Printf("CustomerContractDocumentStore: Create %v\n", err)
		return err
	}

	return nil
}
