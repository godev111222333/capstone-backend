package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CarImageStore struct {
	db *gorm.DB
}

func NewCarImageStore(db *gorm.DB) *CarImageStore {
	return &CarImageStore{db: db}
}

func (s *CarImageStore) Create(images []*model.CarImage) error {
	if err := s.db.Create(images).Error; err != nil {
		fmt.Printf("CarImageStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *CarImageStore) GetByCategory(
	carID int,
	category model.CarImageCategory,
	status model.CarImageStatus,
	limit int,
) ([]string, error) {
	var res []*model.CarImage
	if err := s.db.Where(
		"car_id = ? and category = ? and status = ?",
		carID, string(category), string(status),
	).Order("id desc").Limit(limit).Find(&res).Error; err != nil {
		fmt.Printf("CarImageStore: GetByCategory %v\n", err)
		return nil, err
	}

	n := len(res)
	urls := make([]string, n)
	for index, image := range res {
		urls[n-index-1] = image.URL
	}
	return urls, nil
}
