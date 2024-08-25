package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"time"

	"github.com/godev111222333/capstone-backend/src/api"
	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/service"
	"github.com/godev111222333/capstone-backend/src/store"
	"github.com/redis/go-redis/v9"
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
	newPartnerApprovalQueue := make(chan int, api.ChanBufferSize)
	backgroundJob := service.NewBackgroundService(cfg.BackgroundJob, dbStore, redisClient, newPartnerApprovalQueue)

	go func() {
		if err := backgroundJob.RunPartnerApprovalChecker(); err != nil {
			panic(err)
		}
	}()

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
		newPartnerApprovalQueue,
	)
	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	go func() {
		runCronJob(feCfg, server)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit API server")
	<-stop
}

func runCronJob(feCfg *misc.FEConfig, server *api.Server) {
	fmt.Println("Cron job running ...")
	c := cron.New()
	if _, err := c.AddFunc("0 0 1 * *", func() {
		now := time.Now()
		lastMonth := now.AddDate(0, -1, 0)
		firstDay := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())
		lastDay := firstDay.AddDate(0, 1, 0).Add(-time.Second)

		_ = server.InternalMakeMonthlyPayment(firstDay, lastDay, feCfg.AdminReturnURL)
	}); err != nil {
		panic(err)
	}
}
