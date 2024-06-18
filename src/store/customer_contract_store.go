package store

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CustomerContractStore struct {
	db *gorm.DB
}

func NewCustomerContractStore(db *gorm.DB) *CustomerContractStore {
	return &CustomerContractStore{db: db}
}

func (s *CustomerContractStore) Create(c *model.CustomerContract) error {
	if err := s.db.Create(c).Error; err != nil {
		fmt.Printf("CustomerContractStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *CustomerContractStore) IsOverlap(carID int, desiredStartDate time.Time, desiredEndDate time.Time) (bool, error) {
	// 1. case start_date >= desiredStartDate
	counter := struct {
		Count int `json:"count"`
	}{}
	raw := `select count(*) as count from customer_contracts where car_id = ? and start_date >= ? and ? >= start_date and (status = 'ordered' or status = 'renting' or status = 'completed')`
	if err := s.db.Raw(raw, carID, desiredStartDate, desiredEndDate).Scan(&counter).Error; err != nil {
		fmt.Printf("CustomerContractStore: CheckOverlap case 1 %v\n", err)
		return false, err
	}

	if counter.Count > 0 {
		return true, nil
	}

	// 2. case desiredStartDate >= start_date
	raw = `select count(*) as count from customer_contracts where car_id = ? and ? >= start_date and end_date >= ? and (status = 'ordered' or status = 'renting' or status = 'completed')`
	if err := s.db.Raw(raw, carID, desiredStartDate, desiredStartDate).Scan(&counter).Error; err != nil {
		fmt.Printf("CustomerContractStore: CheckOverlap case 2 %v\n", err)
		return false, err
	}

	return counter.Count > 0, nil
}

func (s *CustomerContractStore) FindByID(id int) (*model.CustomerContract, error) {
	res := &model.CustomerContract{}
	if err := s.db.Where("id = ?", id).Preload("Customer").Preload("Car").Find(res).Error; err != nil {
		fmt.Printf("CustomerContractStore: FindByID %v\n", err)
		return nil, err
	}

	return res, nil
}

func (s *CustomerContractStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(&model.CustomerContract{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("CustomerContractStore: Update %v\n", err)
		return err
	}
	return nil
}
