package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type GarageConfigStore struct {
	db *gorm.DB
}

func NewGarageConfigStore(db *gorm.DB) *GarageConfigStore {
	return &GarageConfigStore{db: db}
}

func (s *GarageConfigStore) Update(configs map[model.GarageConfigType]int) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		for typez, newMax := range configs {
			if newMax == 0 {
				continue
			}

			if err := tx.Model(&model.GarageConfig{}).
				Where("type = ?", typez).
				Update("maximum", newMax).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		fmt.Printf("GarageConfigStore: Update %v\n", err)
		return err
	}

	return nil
}

func (s *GarageConfigStore) Get() (map[model.GarageConfigType]int, error) {
	m := make(map[model.GarageConfigType]int)

	for _, typez := range []model.GarageConfigType{
		model.GarageConfigTypeMax4Seats,
		model.GarageConfigTypeMax7Seats,
		model.GarageConfigTypeMax15Seats,
	} {
		r := &model.GarageConfig{}
		if err := s.db.Where("type = ?", string(typez)).First(r).Error; err != nil {
			return nil, err
		}

		m[typez] = r.Maximum
	}

	return m, nil
}
