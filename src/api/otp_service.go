package api

import (
	"strings"
	"time"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/godev111222333/capstone-backend/src/store"
	"github.com/twilio/twilio-go"
)

const (
	FakeOTP = "999999"
)

type OTPService struct {
	cfg    *misc.OTPConfig
	db     *store.DbStore
	client *twilio.RestClient
}

func NewOTPService(
	cfg *misc.OTPConfig,
	db *store.DbStore,
) *OTPService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		AccountSid: cfg.AccountSID,
		Username:   cfg.ApiKey,
		Password:   cfg.ApiSecret,
	})
	return &OTPService{cfg, db, client}
}

func (s *OTPService) SendOTP(otpType model.OTPType, phoneNumber string) error {
	//phoneWithPrefix := addPhoneCountryPrefix(phoneNumber)
	//code := misc.RandomOTP(6)
	//msgBody := fmt.Sprintf("MinhHungCar verification code: %s. Do not share this code with anyone", code)
	//param := &twilioApi.CreateMessageParams{}
	//param.SetFrom(s.cfg.FromNumber)
	//param.SetTo(phoneWithPrefix)
	//param.SetBody(msgBody)
	//
	//_, err := s.client.Api.CreateMessage(param)
	//if err != nil {
	//	fmt.Printf("OTPService: SentOTP %v\n", err)
	//	return err
	//}

	now := time.Now()
	if err := s.db.OTPStore.Create(&model.OTP{
		OtpType:     otpType,
		PhoneNumber: phoneNumber,
		OTP:         FakeOTP,
		Status:      model.OTPStatusSent,
		ExpiresAt:   now.Add(30 * time.Minute),
		CreatedAt:   now,
		UpdatedAt:   now,
	}); err != nil {
		return err
	}

	return nil
}

func (s *OTPService) VerifyOTP(otpType model.OTPType, phone string, otp string) (bool, error) {
	sentOTP, err := s.db.OTPStore.GetLastByOTPType(phone, otpType)
	if err != nil {
		return false, err
	}

	if otp == sentOTP.OTP && sentOTP.Status == model.OTPStatusSent && sentOTP.ExpiresAt.After(time.Now()) {
		return true, nil
	}

	return false, nil
}

func addPhoneCountryPrefix(phoneNumber string) string {
	if len(phoneNumber) < 1 {
		return ""
	}

	return strings.Replace(phoneNumber, "0", "+84", 1)
}
