package store

import (
	"fmt"
	"testing"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestAccountStore_Get(t *testing.T) {
	accounts := []*model.Account{
		{Email: "acc1@gmail.com", FirstName: "First 1", LastName: "Last 1", PhoneNumber: "phone1", RoleID: model.RoleIDCustomer, Status: model.AccountStatusActive},
		{Email: "acc2@gmail.com", FirstName: "First 2", LastName: "Last 2", PhoneNumber: "phone2", RoleID: model.RoleIDPartner, Status: model.AccountStatusActive},
		{Email: "acc3@gmail.com", FirstName: "First 2", LastName: "Last 3", PhoneNumber: "phone3", RoleID: model.RoleIDCustomer, Status: model.AccountStatusActive},
		{Email: "acc4@gmail.com", FirstName: "First 4", LastName: "Last 4", PhoneNumber: "phone4", RoleID: model.RoleIDCustomer, Status: model.AccountStatusActive},
	}

	for _, acct := range accounts {
		require.NoError(t, TestDb.AccountStore.Create(acct))
	}

	t.Run("list all active accounts", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 3)
	})

	t.Run("list all partners", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, model.RoleNamePartner, "", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 1)
		fmt.Println(accts[0].ID)
	})

	t.Run("list all customers", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, model.RoleNameCustomer, "", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 2)
	})

	t.Run("get account 1", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "phone1", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 1)
	})

	t.Run("get by Last", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "Last", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 3)
	})

	t.Run("get by phone1", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "phone1", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 1)
	})

	t.Run("get by email", func(t *testing.T) {
		accts, err := TestDb.AccountStore.Get(model.AccountStatusActive, "", "acc4@gmail.com", 0, 10)
		require.NoError(t, err)
		require.Len(t, accts, 1)
	})
}
