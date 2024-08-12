package main

import (
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/godev111222333/capstone-backend/src/api"
	"github.com/godev111222333/capstone-backend/src/cmd/seeder"
	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/service"
	"github.com/godev111222333/capstone-backend/src/store"
)

func main() {
	cfg, err := misc.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	feCfg, err := misc.LoadFEConfig("fe-config.yaml")
	if err != nil {
		panic(err)
	}

	dbStore, err := store.NewDbStore(cfg.Database)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	otpService := service.NewOTPService(cfg.OTP, redisClient)
	s3Store := store.NewS3Store(cfg.AWS)
	bankMetadata, err := misc.LoadBankMetadata("etc/converted_banks.txt")
	if err != nil {
		panic(err)
	}

	pdfService := service.NewPDFService(cfg.PDFService)
	paymentService := api.NewVnPayService(cfg.VNPay)
	notificationPushService := service.NewNotificationPushService("", dbStore)

	server := api.NewServer(
		cfg.ApiServer,
		feCfg,
		dbStore,
		s3Store,
		otpService,
		bankMetadata,
		pdfService,
		paymentService,
		notificationPushService,
		redisClient,
		nil,
	)

	if err := seeder.SeedAccounts(dbStore); err != nil {
		panic(err)
	}

	if err := seeder.SeedCustomerContractRules(dbStore); err != nil {
		panic(err)
	}

	if err := seeder.SeedPartnerContractRule(dbStore); err != nil {
		panic(err)
	}

	if err := seeder.SeedCars(server, dbStore); err != nil {
		panic(err)
	}

	if err := seeder.SeedCarImages(dbStore); err != nil {
		panic(err)
	}

	if err := seeder.SeedCustomerContract(server, dbStore); err != nil {
		panic(err)
	}

	if err := seeder.SeedCustomerContractImages(dbStore); err != nil {
		panic(err)
	}

	if err := seeder.SeedCustomerPayments(server, dbStore); err != nil {
		panic(err)
	}

	if err := setSequenceValues(dbStore); err != nil {
		panic(err)
	}

	fmt.Println("Seed data from CSV done!")
}

func setSequenceValues(dbStore *store.DbStore) error {
	sqls := []string{
		`select setval('accounts_id_seq', (select MAX(id) from accounts))`,
		`select setval('cars_id_seq', (select MAX(id) from cars))`,
		`select setval('car_images_id_seq', (select MAX(id) from car_images))`,
		`select setval('customer_contracts_id_seq', (select MAX(id) from customer_contracts))`,
		`select setval('customer_contract_rules_id_seq', (select MAX(id) from customer_contract_rules))`,
		`select setval('partner_contract_rules_id_seq', (select MAX(id) from partner_contract_rules))`,
		`select setval('customer_payments_id_seq', (select MAX(id) from customer_payments))`,
		`select setval('customer_contract_images_id_seq', (select MAX(id) from customer_contract_images))`,
		`update accounts set phone_number = concat('0', phone_number) where phone_number != 'admin'  and phone_number != 'tech'`,
	}
	for _, sql := range sqls {
		if err := dbStore.DB.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}
