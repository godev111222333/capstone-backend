package api

import (
	"testing"

	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/stretchr/testify/require"
)

func TestOTPService(t *testing.T) {
	// Integration test, temporary skip
	t.Skip()
	t.Parallel()

	t.Run("send OTP", func(t *testing.T) {
		t.Parallel()

		phoneNumber := "0389068116"
		require.NoError(t, TestDb.AccountStore.Create(
			&model.Account{
				RoleID:      model.RoleIDPartner,
				PhoneNumber: phoneNumber,
				Status:      model.AccountStatusWaitingConfirmEmail,
			},
		))

		otpService := NewOTPService(TestConfig.OTP, TestDb)
		require.NoError(t, otpService.SendOTP(model.OTPTypeRegister, phoneNumber))
	})
}
