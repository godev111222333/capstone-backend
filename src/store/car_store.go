package store

import (
	"errors"
	"fmt"
	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CarStore struct {
	db *gorm.DB
}

func NewCarStore(db *gorm.DB) *CarStore {
	return &CarStore{db: db}
}

func (s *CarStore) Create(car *model.Car) error {
	if err := s.db.Create(car).Error; err != nil {
		fmt.Printf("CarStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *CarStore) GetByID(id int) (*model.Car, error) {
	res := &model.Car{}
	if err := s.db.Where("id = ?", id).Preload("Account").Find(res).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}
	return res, nil
}

func (s *CarStore) GetByPartner(partnerID, offset, limit int, status model.CarStatus) ([]*model.Car, error) {
	var res []*model.Car
	if limit == 0 {
		limit = 1000
	}

	var row *gorm.DB
	if status == model.CarStatusNoFilter {
		row = s.db.Where("partner_id = ?", partnerID).Preload("Account").Preload("CarModel").Order("id desc").Offset(offset).Limit(limit).Find(&res)
	} else {
		row = s.db.Where("partner_id = ? and status = ?", partnerID, string(status)).Preload("Account").Preload("CarModel").Order("id desc").Offset(offset).Limit(limit).Find(&res)
	}
	if err := row.Error; err != nil {
		fmt.Printf("CarStore: GetByPartner %v\n", err)
		return nil, err
	}

	if row.RowsAffected == 0 {
		return []*model.Car{}, nil
	}

	return res, nil
}

func (s *CarStore) Update(id int, values map[string]interface{}) error {
	if err := s.db.Model(&model.Car{}).Where("id = ?", id).Updates(values).Error; err != nil {
		fmt.Printf("CarStore: Update %v\n", err)
		return err
	}
	return nil
}
