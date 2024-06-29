package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CustomerContractImageStore struct {
	db *gorm.DB
}

func NewCustomerContractImageStore(db *gorm.DB) *CustomerContractImageStore {
	return &CustomerContractImageStore{db: db}
}

func (s *CustomerContractImageStore) Create(images []*model.CustomerContractImage) error {
	if err := s.db.Create(images).Error; err != nil {
		fmt.Printf("CustomerContractImageStore: Create %v\n", err)
		return err
	}

	return nil
}

func (s *CustomerContractImageStore) Get(
	cusContractId int,
	category model.CustomerContractImageCategory,
	limit int,
	status model.CustomerContractImageStatus,
) ([]string, error) {
	res := make([]model.CustomerContractImage, 0)
	if err := s.db.Where(
		"customer_contract_id = ? and category = ? and status = ?",
		cusContractId, string(category), string(status)).Order("id desc").Limit(limit).Scan(&res).Error; err != nil {
		fmt.Printf("CustomerContractImageStore: Get %v\n", err)
		return nil, err
	}

	n := len(res)
	urls := make([]string, n)
	for index, image := range res {
		urls[n-1-index] = image.URL
	}
	return urls, nil
}
