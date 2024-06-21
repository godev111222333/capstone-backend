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

func (s *CustomerPaymentStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(model.CustomerPayment{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("CustomerPaymentStore: Update %v\n", err)
		return err
	}
	return nil
}

func (s *CustomerPaymentStore) GetByID(id int) (*model.CustomerPayment, error) {
	res := &model.CustomerPayment{}
	if err := s.db.Where("id = ?", id).Preload("CustomerContract").Find(res).Error; err != nil {
		fmt.Printf("CustomerPaymentStore: GetByID %v\n", err)
		return nil, err
	}

	return res, nil
}