package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type DrivingLicenseImageStore struct {
	db *gorm.DB
}

func NewDrivingLicenseImageStore(db *gorm.DB) *DrivingLicenseImageStore {
	return &DrivingLicenseImageStore{db: db}
}

func (s *DrivingLicenseImageStore) Create(images []*model.DrivingLicenseImage) error {
	if err := s.db.Create(images).Error; err != nil {
		fmt.Printf("DrivingLicenseImageStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *DrivingLicenseImageStore) Get(
	accID int,
	status model.DrivingLicenseImageStatus,
	limit int,
) ([]*model.DrivingLicenseImage, error) {
	var res []*model.DrivingLicenseImage
	if err := s.db.Where("account_id = ? and status = ?", accID, string(status)).Order("id desc").Limit(limit).Find(&res).Error; err != nil {
		fmt.Printf("DrivingLicenseImageStore: Get %v\n", err)
		return nil, err
	}

	return res, nil
}
