package seeder

import (
	"os"

	"github.com/gocarina/gocsv"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
)

type CustomerContractRule struct {
	ID                   int      `csv:"id"`
	InsurancePercent     float64  `csv:"insurance_percent"`
	PrepayPercent        float64  `csv:"prepay_percent"`
	CollateralCashAmount int      `csv:"collateral_cash_amount"`
	CreatedAt            DateTime `json:"created_at"`
	UpdatedAt            DateTime `json:"updated_at"`
}

func (ccr *CustomerContractRule) ToCustomerContractRuleDB() *model.CustomerContractRule {
	return &model.CustomerContractRule{
		ID:                   ccr.ID,
		InsurancePercent:     ccr.InsurancePercent,
		PrepayPercent:        ccr.PrepayPercent,
		CollateralCashAmount: ccr.CollateralCashAmount,
		CreatedAt:            ccr.CreatedAt.Time,
		UpdatedAt:            ccr.UpdatedAt.Time,
	}
}

func SeedCustomerContractRules(dbStore *store.DbStore) error {
	cusRules := make([]*CustomerContractRule, 0)
	rulesFile, err := os.OpenFile(toFilePath(CustomerContractRulesFile), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer rulesFile.Close()

	if err := gocsv.UnmarshalFile(rulesFile, &cusRules); err != nil {
		return err
	}

	rules := make([]*model.CustomerContractRule, len(cusRules))
	for i, a := range cusRules {
		rules[i] = a.ToCustomerContractRuleDB()
	}

	return dbStore.CustomerContractRuleStore.CreateBatch(rules)
}
