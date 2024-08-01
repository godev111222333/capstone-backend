package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CustomerContractRuleStore struct {
	db *gorm.DB
}

func NewCustomerContractRuleStore(db *gorm.DB) *CustomerContractRuleStore {
	return &CustomerContractRuleStore{db: db}
}

func (s *CustomerContractRuleStore) GetLast() (*model.CustomerContractRule, error) {
	res := &model.CustomerContractRule{}
	if err := s.db.Order("id desc").First(res).Error; err != nil {
		fmt.Printf("CustomerContractRuleStore: GetLast %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *CustomerContractRuleStore) Create(rule *model.CustomerContractRule) error {
	if err := s.db.Create(rule).Error; err != nil {
		fmt.Printf("CustomerContractRuleStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *CustomerContractRuleStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(model.CustomerContractRule{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("CustomerContractRuleStore: Update %v\n", err)
		return err
	}

	return nil
}
