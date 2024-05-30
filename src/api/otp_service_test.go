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

		require.NoError(t, TestDb.PartnerStore.Create(&model.Partner{
			Account: model.Account{
				RoleID: model.RoleIDPartner,
				Email:  "godev111222333@gmail.com",
				Status: model.AccountStatusWaitingConfirmEmail,
			},
		}))

		otpService := NewOTPService(TestDb, TestConfig.OTP.Email, TestConfig.OTP.Password)
		require.NoError(t, otpService.SendOTP(model.OTPTypeRegister, "godev111222333@gmail.com"))
	})
}
