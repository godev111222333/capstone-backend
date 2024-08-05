package seeder

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/godev111222333/capstone-backend/src/api"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type CustomerPayment struct {
	ID                 int                 `csv:"id"`
	CustomerContractID int                 `csv:"customer_contract_id"`
	PaymentType        model.PaymentType   `csv:"payment_type"`
	Amount             int                 `csv:"amount"`
	Note               string              `csv:"note"`
	Status             model.PaymentStatus `csv:"status"`
	PaymentURL         string              `csv:"payment_url"`
	ReturnURL          string              `csv:"return_url"`
	CreatedAt          DateTime            `csv:"created_at"`
	UpdatedAt          DateTime            `csv:"updated_at"`
}

func (cp *CustomerPayment) ToCustomerPaymentDb() *model.CustomerPayment {
	return &model.CustomerPayment{
		ID:                 cp.ID,
		CustomerContractID: cp.CustomerContractID,
		PaymentType:        cp.PaymentType,
		Amount:             cp.Amount,
		Note:               cp.Note,
		Status:             cp.Status,
		PaymentURL:         cp.PaymentURL,
		CreatedAt:          cp.CreatedAt.Time,
		UpdatedAt:          cp.UpdatedAt.Time,
	}
}

func SeedCustomerPayments(server *api.Server, dbStore *store.DbStore) error {
	payments := make([]*CustomerPayment, 0)
	cusPaymentFile, err := os.OpenFile(toFilePath(CustomerPaymentsFile), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer cusPaymentFile.Close()

	if err := gocsv.UnmarshalFile(cusPaymentFile, &payments); err != nil {
		return err
	}

	cusPayments := make([]*model.CustomerPayment, len(payments))
	for i, a := range payments {
		cusPayments[i] = a.ToCustomerPaymentDb()
	}

	if err := dbStore.CustomerPaymentStore.CreateBatch(cusPayments); err != nil {
		return err
	}

	for i, payment := range cusPayments {
		url, err := server.PaymentService.GeneratePaymentURL(
			[]int{payment.ID}, payment.Amount, time.Now().Format("02150405"), payments[i].ReturnURL)
		if err != nil {
			return err
		}

		if err := dbStore.CustomerPaymentStore.Update(payment.ID, map[string]interface{}{"payment_url": url}); err != nil {
			return err
		}
	}

	return nil
}
