package store

import (
	"testing"
	"time"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestPartnerStore(t *testing.T) {
	t.Parallel()

	t.Run("create partner", func(t *testing.T) {
		t.Parallel()

		partner := &model.Partner{
			Account: model.Account{
				RoleID:                   model.RoleIDPartner,
				FirstName:                "Thien",
				LastName:                 "Nguyen",
				PhoneNumber:              "1111111",
				Email:                    "thien@gmail.com",
				IdentificationCardNumber: "2222",
				DateOfBirth:              time.Now(),
				Password:                 "easy_password",
				AvatarURL:                "https://facebook.com",
				Status:                   model.AccountStatusEnable,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		require.NoError(t, TestDb.PartnerStore.Create(partner))

		insertedPartner, err := TestDb.PartnerStore.GetByID(partner.ID)
		require.NoError(t, err)
		require.Equal(t, model.RoleIDPartner, insertedPartner.Account.RoleID)
		require.Equal(t, "Thien", insertedPartner.Account.FirstName)
		require.Equal(t, "Nguyen", insertedPartner.Account.LastName)
		require.Equal(t, "1111111", insertedPartner.Account.PhoneNumber)
	})
}
