package store

import (
	"fmt"
	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CarModelStore struct {
	db *gorm.DB
}

func NewCarModelStore(db *gorm.DB) *CarModelStore {
	return &CarModelStore{db: db}
}

func (s *CarModelStore) Create(models []*model.CarModel) error {
	if err := s.db.Create(models).Error; err != nil {
		fmt.Printf("CarModelStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *CarModelStore) GetAll() ([]*model.CarModel, error) {
	var models []*model.CarModel
	if err := s.db.Find(&models).Error; err != nil {
		fmt.Printf("CarModelStore: SearchCars %v\n", err)
		return nil, err
	}

	return models, nil
}
