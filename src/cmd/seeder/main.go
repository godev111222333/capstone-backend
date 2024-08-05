package main

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/godev111222333/capstone-backend/src/api"
	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/service"
	"github.com/godev111222333/capstone-backend/src/store"
)

const (
	DateTimeLayout = "2006-01-02 15:04:05"
)

type DateTime struct {
	time.Time
}

func (date *DateTime) MarshalCSV() (string, error) {
	return date.Time.Add(7 * time.Hour).Format(DateTimeLayout), nil
}

func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse(DateTimeLayout, csv)
	date.Time = date.Time.Add(-7 * time.Hour)
	return err
}

func toFilePath(file string) string {
	return fmt.Sprintf("etc/seed/%s", file)
}

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

	if err := seedAccounts(dbStore); err != nil {
		panic(err)
	}

	if err := seedCars(server, dbStore); err != nil {
		panic(err)
	}

	if err := seedCarImages(dbStore); err != nil {
		panic(err)
	}

	if err := seedCustomerContract(server, dbStore); err != nil {
		panic(err)
	}

	if err := seedCustomerContractImages(dbStore); err != nil {
		panic(err)
	}

	if err := seedCustomerPayments(server, dbStore); err != nil {
		panic(err)
	}

	fmt.Println("Seed data from CSV done!")
}
