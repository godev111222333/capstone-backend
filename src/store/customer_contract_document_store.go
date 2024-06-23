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

func (s *CustomerContractDocumentStore) GetByCategory(
	cusContractID int,
	category model.DocumentCategory,
	limit int,
	status model.DocumentStatus,
) ([]*model.Document, error) {
	var res []*model.Document
	rawSql := `select * from documents where category = ? and status = ? and id in (select document_id from customer_contract_documents where customer_contract_id = ?) limit ?`
	if err := s.db.Raw(rawSql, string(category), string(status), cusContractID, limit).Scan(&res).Error; err != nil {
		fmt.Printf("CustomerContractDocumentStore: GetByCategory %v\n", err)
		return nil, err
	}

	return res, nil
}
