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

	fmt.Println("Seed data from CSV done!")
}
