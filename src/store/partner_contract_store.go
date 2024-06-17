package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type PartnerContractStore struct {
	db *gorm.DB
}

func NewPartnerContractStore(db *gorm.DB) *PartnerContractStore {
	return &PartnerContractStore{db: db}
}

func (s *PartnerContractStore) Create(c *model.PartnerContract) error {
	if err := s.db.Create(c).Error; err != nil {
		fmt.Printf("PartnerContractStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *PartnerContractStore) GetByCarID(carID int) (*model.PartnerContract, error) {
	res := &model.PartnerContract{}
	if err := s.db.Where("car_id = ?", carID).Preload("Car").Preload("Car.CarModel").Find(res).Error; err != nil {
		fmt.Printf("PartnerContractStore: GetByCarID %v\n", err)
		return nil, err
	}
	return res, nil
}

func (s *PartnerContractStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(model.PartnerContract{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("PartnerContractStore: Update %v\n", err)
		return err
	}
	return nil
}
