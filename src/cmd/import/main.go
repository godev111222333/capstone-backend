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

	models, err := importFromFile("etc/car_models/2024.csv")
	if err != nil {
		panic(err)
	}

	if err := dbStore.CarModelStore.Create(models); err != nil {
		panic(err)
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
		})
	}

	return models, nil
}

func randomSeats() int {
	seats := []int{4, 7, 15}
	rand.Seed(time.Now().UnixNano())
	return seats[rand.Intn(len(seats))]
}
