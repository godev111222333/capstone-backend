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
) ([]*model.CustomerContractImage, error) {
	res := make([]*model.CustomerContractImage, 0)
	if err := s.db.Where(
		"customer_contract_id = ? and category = ? and status = ?",
		cusContractId, string(category), string(status)).Order("id desc").Limit(limit).Find(&res).Error; err != nil {
		fmt.Printf("CustomerContractImageStore: Get %v\n", err)
		return nil, err
	}

	n := len(res)
	respImages := make([]*model.CustomerContractImage, n)
	for index, image := range res {
		respImages[n-1-index] = image
	}
	return respImages, nil
}

func (s *CustomerContractImageStore) Update(imageID int, newStatus model.CustomerContractImageStatus) error {
	if err := s.db.Model(model.CustomerContractImage{}).
		Where("id = ?", imageID).
		Updates(map[string]interface{}{"status": string(newStatus)}).Error; err != nil {
		fmt.Printf("CustomerContractImageStore: Update %v\n", err)
		return err
	}

	return nil
}
