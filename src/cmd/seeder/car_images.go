package main

import (
	"os"

	"github.com/gocarina/gocsv"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type CarImage struct {
	ID        int                    `csv:"id"`
	CarID     int                    `csv:"car_id"`
	Car       *Car                   `csv:"car,omitempty"`
	URL       string                 `csv:"url"`
	Category  model.CarImageCategory `csv:"category"`
	Status    model.CarImageStatus   `csv:"status"`
	CreatedAt DateTime               `csv:"created_at"`
	UpdatedAt DateTime               `csv:"updated_at"`
}

func (ci *CarImage) ToDbCarImage() *model.CarImage {
	return &model.CarImage{
		ID:        ci.ID,
		CarID:     ci.CarID,
		URL:       ci.URL,
		Category:  ci.Category,
		Status:    ci.Status,
		CreatedAt: ci.CreatedAt.Time,
		UpdatedAt: ci.UpdatedAt.Time,
	}
}

func seedCarImages(dbStore *store.DbStore) error {
	images := make([]*CarImage, 0)
	carImageFile, err := os.OpenFile(toFilePath(CarImagesFile), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer carImageFile.Close()

	if err := gocsv.UnmarshalFile(carImageFile, &images); err != nil {
		return err
	}

	carImages := make([]*model.CarImage, len(images))
	for i, a := range images {
		carImages[i] = a.ToDbCarImage()
	}

	return dbStore.CarImageStore.Create(carImages)
}
