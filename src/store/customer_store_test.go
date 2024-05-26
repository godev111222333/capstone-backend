package store

import (
	"testing"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestCustomerStore(t *testing.T) {
	t.Parallel()

	t.Run("create customer", func(t *testing.T) {
		t.Parallel()

		customer := &model.Customer{
			Account: model.Account{
				RoleID:                   model.RoleIDCustomer,
				FirstName:                "Son",
				LastName:                 "Le",
				PhoneNumber:              "0987654321",
				Email:                    "son@gmail.com",
				IdentificationCardNumber: "12345",
				Password:                 "hard_password",
				AvatarURL:                "https://google.com",
				Status:                   model.AccountStatusEnable,
			},
			DrivingLicense: "9876",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		require.NoError(t, TestDb.CustomerStore.Create(customer))

		insertedCus, err := TestDb.CustomerStore.GetByID(customer.ID)
		require.NoError(t, err)
		require.Equal(t, "Son", insertedCus.Account.FirstName)
		require.Equal(t, "Le", insertedCus.Account.LastName)
		require.Equal(t, "0987654321", insertedCus.Account.PhoneNumber)
		require.Equal(t, "9876", insertedCus.DrivingLicense)
	})
}
