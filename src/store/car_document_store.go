package store

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/model"
)

type CarDocumentStore struct {
	db *gorm.DB
}

func NewCarDocumentStore(db *gorm.DB) *CarDocumentStore {
	return &CarDocumentStore{db: db}
}

func (s *CarDocumentStore) Create(carID int, document *model.Document) error {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(document).Error; err != nil {
			return err
		}
		return tx.Create(&model.CarDocument{
			CarID:      carID,
			DocumentID: document.ID,
		}).Error
	}); err != nil {
		fmt.Printf("CarDocumentStore: Create %v\n", err)
		return err
	}
	return nil
}

func (s *CarDocumentStore) GetThumbnailURL(carID int) (string, error) {
	rawSQL := `select * from documents where status = ? and category = ? and id in (select id from car_documents where car_id = ?) order by id asc limit 1;`
	res := &model.Document{}
	if err := s.db.Raw(rawSQL, model.DocumentStatusActive, model.DocumentCategoryCarImages, carID).Scan(res).Error; err != nil {
		fmt.Printf("CarDocumentStore: GetThumbnail %v\n", err)
		return "", err
	}

	return res.Url, nil
}
