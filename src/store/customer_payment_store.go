package store

import (
	"fmt"

	"github.com/godev111222333/capstone-backend/src/model"
	"gorm.io/gorm"
)

type CustomerPaymentStore struct {
	db *gorm.DB
}

func NewCustomerPaymentStore(db *gorm.DB) *CustomerPaymentStore {
	return &CustomerPaymentStore{db: db}
}

func (s *CustomerPaymentStore) Create(m *model.CustomerPayment) error {
	if err := s.db.Create(m).Error; err != nil {
		fmt.Printf("CustomerPaymentStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *CustomerPaymentStore) CreatePaymentDocument(customerPaymentID int, docID int) error {
	m := &model.CustomerPaymentDocument{
		CustomerPaymentID: customerPaymentID,
		DocumentID:        docID,
	}
	if err := s.db.Create(m).Error; err != nil {
		fmt.Printf("CustomerContractStore: CreatePaymentDocument %v\n", err)
		return err
	}

	return nil
}
