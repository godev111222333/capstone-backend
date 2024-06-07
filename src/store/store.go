package store

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/misc"
)

type DbStore struct {
	DB                      *gorm.DB
	AccountStore            *AccountStore
	OTPStore                *OTPStore
	SessionStore            *SessionStore
	CarModelStore           *CarModelStore
	CarStore                *CarStore
	CarDocumentStore        *CarDocumentStore
	PaymentInformationStore *PaymentInformationStore
	GarageConfigStore       *GarageConfigStore
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
		DB:                      db,
		AccountStore:            NewAccountStore(db),
		OTPStore:                NewOTPStore(db),
		SessionStore:            NewSessionStore(db),
		CarModelStore:           NewCarModelStore(db),
		CarStore:                NewCarStore(db),
		CarDocumentStore:        NewCarDocumentStore(db),
		PaymentInformationStore: NewPaymentInformationStore(db),
		GarageConfigStore:       NewGarageConfigStore(db),
	}, nil
}
