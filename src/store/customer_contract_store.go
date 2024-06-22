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

func (s *CustomerContractStore) GetByCustomerID(cusID int, status model.CustomerContractStatus, offset, limit int) ([]*model.CustomerContract, error) {
	var res []*model.CustomerContract
	if limit == 0 {
		limit = 1000
	}

	if status == model.CustomerContractStatusNoFilter {
		if err := s.db.Where("customer_id = ?", cusID).Preload("Customer").Preload("Car").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerContractStore: GetByCustomerID %v\n", err)
			return nil, err
		}
	} else {
		if err := s.db.Where("customer_id = ? and status like ?", cusID, "%"+string(status)+"%").Preload("Customer").Preload("Car").Offset(offset).Limit(limit).Find(&res).Error; err != nil {
			fmt.Printf("CustomerContractStore: GetByCustomerID %v\n", err)
			return nil, err
		}
	}

	return res, nil
}

func (s *CustomerContractStore) GetByStatus(status model.CustomerContractStatus, offset, limit int) ([]*model.CustomerContract, error) {
	res := []*model.CustomerContract{}
	var err error
	if limit == 0 {
		limit = 1000
	}

	if status == model.CustomerContractStatusNoFilter {
		err = s.db.Preload("Customer").Preload("Car").Offset(offset).Limit(limit).Find(&res).Error
	} else {
		err = s.db.Where("status = ?", string(status)).Preload("Customer").Preload("Car").Offset(offset).Limit(limit).Find(&res).Error
	}

	if err != nil {
		fmt.Printf("CustomerContractStore: GetByStatus %v\n", err)
		return nil, err
	}

	return res, nil
}
