package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

func main() {
	cfg, err := misc.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}
	dbStore, err := store.NewDbStore(cfg.Database)
	if err != nil {
		panic(err)
	}

	filePaths := make([]string, 0)
	for year := 2010; year <= 2024; year++ {
		filePaths = append(filePaths, fmt.Sprintf("etc/car_models/%d.csv", year))
	}
	for _, path := range filePaths {
		models, err := importFromFile(path)
		if err != nil {
			panic(err)
		}

		if err := dbStore.CarModelStore.Create(models); err != nil {
			panic(err)
		}
	}

	fmt.Println("imported car models data successfully")
}

func importFromFile(filePath string) ([]*model.CarModel, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	// Convert CSV records to Account structs
	var models []*model.CarModel
	for line, rc := range records {
		if line == 0 {
			continue
		}
		year, err := strconv.Atoi(rc[0])
		if err != nil {
			fmt.Println("Error converting year:", err)
			continue
		}
		models = append(models, &model.CarModel{
			Brand:         rc[1],
			Model:         rc[2],
			Year:          year,
			NumberOfSeats: randomSeats(),
			BasedPrice:    randomBasedPrice(),
		})
	}

	return models, nil
}

func randomSeats() int {
	seats := []int{4, 7, 15}
	rand.Seed(time.Now().UnixNano())
	return seats[rand.Intn(len(seats))]
}

// random price: 200k -> 1500k
func randomBasedPrice() int {
	return (200 + rand.Intn(1300)) * 1000
}
