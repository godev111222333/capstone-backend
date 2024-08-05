package seeder

import (
	"os"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/godev111222333/capstone-backend/src/api"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type CustomerContract struct {
	ID                      int                          `csv:"id"`
	CustomerID              int                          `csv:"customer_id"`
	CarID                   int                          `csv:"car_id"`
	StartDate               DateTime                     `csv:"start_date"`
	EndDate                 DateTime                     `csv:"end_date"`
	Status                  model.CustomerContractStatus `csv:"status"`
	Reason                  string                       `csv:"reason"`
	RentPrice               int                          `csv:"rent_price"`
	InsuranceAmount         int                          `csv:"insurance_amount"`
	CollateralType          model.CollateralType         `csv:"collateral_type"`
	IsReturnCollateralAsset bool                         `csv:"is_return_collateral_asset"`
	Url                     string                       `csv:"url"`
	BankName                string                       `csv:"bank_name"`
	BankNumber              string                       `csv:"bank_number"`
	BankOwner               string                       `csv:"bank_owner"`
	CustomerContractRuleID  int                          `csv:"customer_contract_rule_id"`
	FeedbackContent         string                       `csv:"feedback_content"`
	FeedbackRating          int                          `csv:"feedback_rating"`
	FeedbackStatus          model.FeedBackStatus         `csv:"feedback_status"`
	CreatedAt               DateTime                     `csv:"created_at"`
	UpdatedAt               DateTime                     `csv:"updated_at"`
}

func (cc *CustomerContract) ToDbCustomerContract() *model.CustomerContract {
	return &model.CustomerContract{
		ID:                      cc.ID,
		CustomerID:              cc.CustomerID,
		CarID:                   cc.CarID,
		StartDate:               cc.StartDate.Time,
		EndDate:                 cc.EndDate.Time,
		Status:                  cc.Status,
		Reason:                  cc.Reason,
		RentPrice:               cc.RentPrice,
		InsuranceAmount:         cc.InsuranceAmount,
		CollateralType:          cc.CollateralType,
		IsReturnCollateralAsset: cc.IsReturnCollateralAsset,
		Url:                     cc.Url,
		BankName:                cc.BankName,
		BankNumber:              cc.BankNumber,
		BankOwner:               cc.BankOwner,
		CustomerContractRuleID:  cc.CustomerContractRuleID,
		FeedbackContent:         cc.FeedbackContent,
		FeedbackRating:          cc.FeedbackRating,
		FeedbackStatus:          cc.FeedbackStatus,
		CreatedAt:               cc.CreatedAt.Time,
		UpdatedAt:               cc.UpdatedAt.Time,
	}
}

func SeedCustomerContract(server *api.Server, dbStore *store.DbStore) error {
	customerContracts := make([]*CustomerContract, 0)
	customerContractFile, err := os.OpenFile(toFilePath(CustomerContractsFile), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer customerContractFile.Close()

	if err := gocsv.UnmarshalFile(customerContractFile, &customerContracts); err != nil {
		return err
	}

	customerContractsDb := make([]*model.CustomerContract, len(customerContracts))
	for i, a := range customerContracts {
		customerContractsDb[i] = a.ToDbCustomerContract()
	}

	if err := dbStore.CustomerContractStore.CreateBatch(customerContractsDb); err != nil {
		return err
	}

	for _, cc := range customerContractsDb {
		contract, err := dbStore.CustomerContractStore.FindByID(cc.ID)
		if err != nil {
			return err
		}

		now := contract.CreatedAt.Add(time.Hour)
		if err := server.InternalRenderCustomerContractPDF(contract.Customer, &contract.Car, contract, now); err != nil {
			return err
		}
	}

	return nil
}
