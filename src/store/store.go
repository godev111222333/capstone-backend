package store

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/godev111222333/capstone-backend/src/misc"
)

type DbStore struct {
	DB                         *gorm.DB
	AccountStore               *AccountStore
	CarModelStore              *CarModelStore
	CarStore                   *CarStore
	CarImageStore              *CarImageStore
	GarageConfigStore          *GarageConfigStore
	PartnerContractStore       *PartnerContractStore
	CustomerContractStore      *CustomerContractStore
	CustomerContractImageStore *CustomerContractImageStore
	CustomerPaymentStore       *CustomerPaymentStore
	DrivingLicenseImageStore   *DrivingLicenseImageStore
	ContractRuleStore          *ContractRuleStore
	ConversationStore          *ConversationStore
	MessageStore               *MessageStore
	NotificationStore          *NotificationStore
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
		DB:                         db,
		AccountStore:               NewAccountStore(db),
		CarModelStore:              NewCarModelStore(db),
		CarStore:                   NewCarStore(db),
		CarImageStore:              NewCarImageStore(db),
		GarageConfigStore:          NewGarageConfigStore(db),
		PartnerContractStore:       NewPartnerContractStore(db),
		CustomerContractStore:      NewCustomerContractStore(db),
		CustomerContractImageStore: NewCustomerContractImageStore(db),
		CustomerPaymentStore:       NewCustomerPaymentStore(db),
		DrivingLicenseImageStore:   NewDrivingLicenseImageStore(db),
		ContractRuleStore:          NewContractRuleStore(db),
		ConversationStore:          NewConversationStore(db),
		MessageStore:               NewMessageStore(db),
		NotificationStore:          NewNotificationStore(db),
	}, nil
}
