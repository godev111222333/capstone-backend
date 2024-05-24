package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CustomerStore struct {
	db *gorm.DB
}

func NewCustomerStore(db *gorm.DB) *CustomerStore {
	return &CustomerStore{db: db}
}

func (s *CustomerStore) Create(cus *model.Customer) error {
	if err := s.db.Create(cus).Error; err != nil {
		fmt.Printf("CustomerStore: %v\n", err)
		return err
	}

	return nil
}

func (s *CustomerStore) GetByID(cusID int) (*model.Customer, error) {
	res := &model.Customer{}
	if err := s.db.Where("id = ?", cusID).Preload("Account").Find(res).Error; err != nil {
		fmt.Printf("CustomerStore: %v\n", err)
		return nil, err
	}

	return res, nil
}
