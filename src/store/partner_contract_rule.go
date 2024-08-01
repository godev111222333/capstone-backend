package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type PartnerContractRuleStore struct {
	db *gorm.DB
}

func NewPartnerContractRuleStore(db *gorm.DB) *PartnerContractRuleStore {
	return &PartnerContractRuleStore{db: db}
}

func (s *PartnerContractRuleStore) GetLast() (*model.PartnerContractRule, error) {
	res := &model.PartnerContractRule{}
	if err := s.db.Order("id desc").First(res).Error; err != nil {
		fmt.Printf("PartnerContractRuleStore: GetLast %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *PartnerContractRuleStore) Create(rule *model.PartnerContractRule) error {
	if err := s.db.Create(rule).Error; err != nil {
		fmt.Printf("PartnerContractRuleStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *PartnerContractRuleStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(model.PartnerContractRule{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("PartnerContractRuleStore: Update %v\n", err)
		return err
	}

	return nil
}
