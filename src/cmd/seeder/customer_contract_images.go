package main

import (
	"os"

	"github.com/gocarina/gocsv"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type CustomerContractImage struct {
	ID                 int                                 `csv:"id"`
	CustomerContractID int                                 `csv:"customer_contract_id"`
	URL                string                              `csv:"url"`
	Category           model.CustomerContractImageCategory `csv:"category"`
	Status             model.CustomerContractImageStatus   `csv:"status"`
	CreatedAt          DateTime                            `csv:"created_at"`
	UpdatedAt          DateTime                            `csv:"updated_at"`
}

func (cci *CustomerContractImage) toCustomerContractImageDB() *model.CustomerContractImage {
	return &model.CustomerContractImage{
		ID:                 cci.ID,
		CustomerContractID: cci.CustomerContractID,
		URL:                cci.URL,
		Category:           cci.Category,
		Status:             cci.Status,
		CreatedAt:          cci.CreatedAt.Time,
		UpdatedAt:          cci.UpdatedAt.Time,
	}
}

func seedCustomerContractImages(dbStore *store.DbStore) error {
	images := make([]*CustomerContractImage, 0)
	customerContractImageFile, err := os.OpenFile(toFilePath(CustomerContractImagesFile), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer customerContractImageFile.Close()

	if err := gocsv.UnmarshalFile(customerContractImageFile, &images); err != nil {
		return err
	}

	cusImages := make([]*model.CustomerContractImage, len(images))
	for i, a := range images {
		cusImages[i] = a.toCustomerContractImageDB()
	}

	return dbStore.CustomerContractImageStore.Create(cusImages)
}
