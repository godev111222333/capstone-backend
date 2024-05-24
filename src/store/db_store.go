package store

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/misc"
)

type DbStore struct {
	DB            *gorm.DB
	CustomerStore *CustomerStore
	PartnerStore  *PartnerStore
}

func NewDbStore(cfg *misc.DatabaseConfig) (*DbStore, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.DbHost, cfg.DbUsername, cfg.DbPassword, cfg.DbName, cfg.DbPort),
	}), &gorm.Config{})

	if err != nil {
		fmt.Printf("error when initing connect to DB: %v\n", err)
		return nil, err
	}

	return &DbStore{
		DB:            db,
		CustomerStore: NewCustomerStore(db),
		PartnerStore:  NewPartnerStore(db),
	}, nil
}
