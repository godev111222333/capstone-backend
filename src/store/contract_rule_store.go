package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type ContractRuleStore struct {
	db *gorm.DB
}

func NewContractRuleStore(db *gorm.DB) *ContractRuleStore {
	return &ContractRuleStore{db: db}
}

func (s *ContractRuleStore) GetLast() (*model.ContractRule, error) {
	res := &model.ContractRule{}
	if err := s.db.Order("id desc").First(res).Error; err != nil {
		fmt.Printf("ContractRuleStore: GetLast %v\n", err)
		return nil, err
	}

	return res, nil
}
