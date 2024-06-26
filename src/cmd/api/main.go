package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/godev111222333/capstone-backend/src/api"
	"github.com/godev111222333/capstone-backend/src/misc"
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

	otpService := api.NewOTPService(cfg.OTP, dbStore)
	s3Store := store.NewS3Store(cfg.AWS)
	bankMetadata, err := misc.LoadBankMetadata("etc/converted_banks.txt")
	if err != nil {
		panic(err)
	}

	pdfService := api.NewPDFService(cfg.PDFService)
	paymentService := api.NewVnPayService(cfg.VNPay)

	server := api.NewServer(cfg.ApiServer, dbStore, s3Store, otpService, bankMetadata, pdfService, paymentService)
	go func() {
		if err := server.Run(); err != nil {
			panic(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit API server")
	<-stop
}
