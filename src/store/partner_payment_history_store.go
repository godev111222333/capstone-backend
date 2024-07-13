package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type PartnerPaymentHistoryStore struct {
	db *gorm.DB
}

func NewPartnerPaymentHistoryStore(db *gorm.DB) *PartnerPaymentHistoryStore {
	return &PartnerPaymentHistoryStore{db: db}
}

func (s *PartnerPaymentHistoryStore) Create(history *model.PartnerPaymentHistory, cusContractIds []int) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(history).Error; err != nil {
			return err
		}

		for _, cusContractID := range cusContractIds {
			m := &model.PartnerPaymentCustomerContract{
				PartnerPaymentHistoryID: history.ID,
				CustomerContractID:      cusContractID,
			}

			if err := tx.Create(m).Error; err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fmt.Printf("PartnerPaymentHistoryStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *PartnerPaymentHistoryStore) GetByID(id int) (*model.PartnerPaymentHistory, error) {
	var res *model.PartnerPaymentHistory
	if err := s.db.Where("id = ?", id).First(res).Error; err != nil {
		fmt.Printf("PartnerPaymentHistoryStore: GetByID %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *PartnerPaymentHistoryStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(model.PartnerPaymentHistory{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("PartnerPaymentHistoryStore: Update %v\n", err)
		return err
	}

	return nil
}

func (s *PartnerPaymentHistoryStore) UpdateMulti(ids []int, values map[string]interface{}) error {
	if err := s.db.Model(model.PartnerPaymentHistory{}).Where("id in ?", ids).Updates(values).Error; err != nil {
		fmt.Printf("PartnerPaymentHistoryStore: UpdateMulti %v\n", err)
		return err
	}

	return nil
}