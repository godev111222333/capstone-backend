package seeder

import (
	"github.com/gocarina/gocsv"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
	"os"
)

type PartnerContractRule struct {
	ID                    int      `csv:"id"`
	RevenueSharingPercent float64  `csv:"revenue_sharing_percent"`
	MaxWarningCount       int      `csv:"max_warning_count"`
	CreatedAt             DateTime `csv:"created_at"`
	UpdatedAt             DateTime `csv:"updated_at"`
}

func (pcr *PartnerContractRule) ToPartnerContractRuleDB() *model.PartnerContractRule {
	return &model.PartnerContractRule{
		ID:                    pcr.ID,
		RevenueSharingPercent: pcr.RevenueSharingPercent,
		MaxWarningCount:       pcr.MaxWarningCount,
		CreatedAt:             pcr.CreatedAt.Time,
		UpdatedAt:             pcr.UpdatedAt.Time,
	}
}

func SeedPartnerContractRule(dbStore *store.DbStore) error {
	partnerRules := make([]*PartnerContractRule, 0)
	rulesFile, err := os.OpenFile(toFilePath(PartnerContractRulesFile), os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer rulesFile.Close()

	if err := gocsv.UnmarshalFile(rulesFile, &partnerRules); err != nil {
		return err
	}

	rules := make([]*model.PartnerContractRule, len(partnerRules))
	for i, a := range partnerRules {
		rules[i] = a.ToPartnerContractRuleDB()
	}

	return dbStore.PartnerContractRuleStore.CreateBatch(rules)
}
