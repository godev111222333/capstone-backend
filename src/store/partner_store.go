package store

import (
	"fmt"
	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type PartnerStore struct {
	db *gorm.DB
}

func NewPartnerStore(db *gorm.DB) *PartnerStore {
	return &PartnerStore{db}
}

func (s *PartnerStore) Create(p *model.Partner) error {
	if err := s.db.Create(p).Error; err != nil {
		fmt.Printf("PartnerStore: %v\n", err)
		return err
	}

	return nil
}

func (s *PartnerStore) GetByID(partnerID int) (*model.Partner, error) {
	res := &model.Partner{}
	if err := s.db.Where("id = ?", partnerID).Preload("Account").Find(res).Error; err != nil {
		fmt.Printf("PartnerStore: %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *PartnerStore) Update(accountID int, values map[string]interface{}) error {
	if err := s.db.Model(&model.Account{ID: accountID}).Updates(values).Error; err != nil {
		fmt.Printf("PartnerStore: %v\n", err)
		return err
	}

	return nil
}
