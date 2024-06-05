package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type PaymentInformationStore struct {
	db *gorm.DB
}

func NewPaymentInformationStore(db *gorm.DB) *PaymentInformationStore {
	return &PaymentInformationStore{db: db}
}

func (s *PaymentInformationStore) Create(p *model.PaymentInformation) error {
	if err := s.db.Create(p).Error; err != nil {
		fmt.Printf("PaymentInformationStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *PaymentInformationStore) Update(acctID int, values map[string]interface{}) error {
	if err := s.db.Model(&model.PaymentInformation{}).Where("account_id = ?", acctID).Updates(values).Error; err != nil {
		fmt.Printf("PaymentInformationStore: Update %v\n", err)
		return err
	}

	return nil
}

func (s *PaymentInformationStore) GetByAcctID(acctID int) (*model.PaymentInformation, error) {
	res := &model.PaymentInformation{}
	if err := s.db.Where("account_id = ?", acctID).First(res).Error; err != nil {
		fmt.Printf("PaymentInformationStore: GetByAcctID %v\n", err)
		return nil, err
	}

	return res, nil
}
