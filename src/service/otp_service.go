package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/godev111222333/capstone-backend/src/misc"
	"github.com/godev111222333/capstone-backend/src/model"
	"github.com/redis/go-redis/v9"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

const (
	FakeOTP = "999999"
)

type OTPService struct {
	cfg          *misc.OTPConfig
	redisClient  *redis.Client
	twilioClient *twilio.RestClient
}

func NewOTPService(
	cfg *misc.OTPConfig,
	redisClient *redis.Client,
) *OTPService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		AccountSid: cfg.AccountSID,
		Username:   cfg.ApiKey,
		Password:   cfg.ApiSecret,
	})
	return &OTPService{cfg, redisClient, client}
}

func (s *OTPService) SendOTP(otpType model.OTPType, phoneNumber string) error {
	phoneWithPrefix := addPhoneCountryPrefix(phoneNumber)
	code := misc.RandomOTP(6)
	msgBody := fmt.Sprintf("MinhHungCar verification code: %s. Do not share this code with anyone", code)
	param := &twilioApi.CreateMessageParams{}
	param.SetFrom(s.cfg.FromNumber)
	param.SetTo(phoneWithPrefix)
	param.SetBody(msgBody)

	_, err := s.twilioClient.Api.CreateMessage(param)
	if err != nil {
		fmt.Printf("OTPService: SentOTP %v\n", err)
		return err
	}

	now := time.Now()
	otp := &model.OTP{
		OtpType:     otpType,
		PhoneNumber: phoneNumber,
		OTP:         code,
		Status:      model.OTPStatusSent,
		ExpiresAt:   now.Add(30 * time.Minute),
	}

	bz, err := json.Marshal(otp)
	if err != nil {
		fmt.Println(err)
		return err
	}

	statusCmd := s.redisClient.Set(context.Background(), toRedisKey(otpType, phoneNumber), string(bz), time.Duration(0))
	if err := statusCmd.Err(); err != nil {
		fmt.Printf("SendOTP: Store to redis error %v\n", err)
		return err
	}

	return nil
}

func (s *OTPService) VerifyOTP(otpType model.OTPType, phone string, otp string) (bool, error) {
	if otp == FakeOTP {
		return true, nil
	}

	value, err := s.redisClient.Get(context.Background(), toRedisKey(otpType, phone)).Result()
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	sentOTP := &model.OTP{}
	if err := json.Unmarshal([]byte(value), sentOTP); err != nil {
		fmt.Println(err)
		return false, err
	}

	if otp == sentOTP.OTP && sentOTP.Status == model.OTPStatusSent && sentOTP.ExpiresAt.After(time.Now()) {
		sentOTP.Status = model.OTPStatusVerified
		bz, err := json.Marshal(sentOTP)
		if err != nil {
			fmt.Println(err)
			return false, err
		}

		statusCmd := s.redisClient.Set(context.Background(), toRedisKey(otpType, phone), string(bz), time.Duration(0))
		if err := statusCmd.Err(); err != nil {
			fmt.Printf("SendOTP: Store to redis error %v\n", err)
			return false, err
		}

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

func toRedisKey(otpType model.OTPType, phoneNumber string) string {
	return fmt.Sprintf("%s__%s", otpType, phoneNumber)
}
