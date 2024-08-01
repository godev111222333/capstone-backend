package store

import (
	"fmt"
	"testing"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestCustomerContractStore(t *testing.T) {
	carModel := &model.CarModel{Brand: "AudiVIP"}
	require.NoError(t, TestDb.CarModelStore.Create([]*model.CarModel{carModel}))
	partner := &model.Account{PhoneNumber: "0001", Status: model.AccountStatusActive, RoleID: model.RoleIDPartner}
	require.NoError(t, TestDb.AccountStore.Create(partner))
	customer := &model.Account{PhoneNumber: "0002", Status: model.AccountStatusActive, RoleID: model.RoleIDCustomer}
	require.NoError(t, TestDb.AccountStore.Create(customer))
	car := &model.Car{PartnerID: partner.ID, CarModelID: carModel.ID, LicensePlate: "8xxx", Status: model.CarStatusActive}
	require.NoError(t, TestDb.CarStore.Create(car))
	contractRuleID := 1

	// order from 10h -> 13h
	now := time.Now()
	contract := &model.CustomerContract{
		CustomerID:     customer.ID,
		CarID:          car.ID,
		StartDate:      now.Add(10 * time.Hour),
		EndDate:        now.Add(13 * time.Hour),
		Status:         model.CustomerContractStatusOrdered,
		ContractRuleID: contractRuleID,
	}
	require.NoError(t, TestDb.CustomerContractStore.Create(contract))

	// order from 18 -> 22h
	contract = &model.CustomerContract{
		CustomerID:     customer.ID,
		CarID:          car.ID,
		StartDate:      now.Add(18 * time.Hour),
		EndDate:        now.Add(22 * time.Hour),
		Status:         model.CustomerContractStatusOrdered,
		ContractRuleID: contractRuleID,
	}
	require.NoError(t, TestDb.CustomerContractStore.Create(contract))

	testCases := []struct {
		desiredStartDate int
		desiredEndDate   int
		isOverlap        bool
	}{
		{
			desiredStartDate: 0,
			desiredEndDate:   1,
			isOverlap:        false,
		},
		{
			desiredStartDate: 1,
			desiredEndDate:   12,
			isOverlap:        true,
		},
		{
			desiredStartDate: 9,
			desiredEndDate:   11,
			isOverlap:        true,
		},
		{
			desiredStartDate: 11,
			desiredEndDate:   14,
			isOverlap:        true,
		},
		{
			desiredStartDate: 13,
			desiredEndDate:   14,
			isOverlap:        true,
		},
		{
			desiredStartDate: 14,
			desiredEndDate:   17,
			isOverlap:        false,
		},
		{
			desiredStartDate: 17,
			desiredEndDate:   23,
			isOverlap:        true,
		},
		{
			desiredStartDate: 20,
			desiredEndDate:   25,
			isOverlap:        true,
		},
		{
			desiredStartDate: 24,
			desiredEndDate:   25,
			isOverlap:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("from %d to %d", tc.desiredStartDate, tc.desiredEndDate), func(t *testing.T) {
			isOverlap, err := TestDb.CustomerContractStore.IsOverlap(
				car.ID,
				now.Add(time.Hour*time.Duration(tc.desiredStartDate)),
				now.Add(time.Hour*time.Duration(tc.desiredEndDate)),
			)
			require.NoError(t, err)
			require.Equal(t, tc.isOverlap, isOverlap)
		})
	}
}
