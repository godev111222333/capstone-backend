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
	if err := s.db.Order("id asc").Find(&models).Error; err != nil {
		fmt.Printf("CarModelStore: GetAll %v\n", err)
		return nil, err
	}

	return models, nil
}

func (s *CarModelStore) GetPagination(offset, limit int) ([]*model.CarModel, error) {
	var models []*model.CarModel
	if err := s.db.Order("id desc").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		fmt.Printf("CarModelStore: GetPagination %v\n", err)
		return nil, err
	}

	return models, nil
}

func (s *CarModelStore) CountTotal() (int, error) {
	var count int64
	if err := s.db.Model(model.CarModel{}).Count(&count).Error; err != nil {
		fmt.Printf("CarModelStore: CountTotal: %v\n", err)
		return 01, err
	}

	return int(count), nil
}

func (s *CarModelStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(model.CarModel{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("CarModelStore: Update %v\n", err)
		return err
	}

	return nil
}
