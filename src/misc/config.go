package misc

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type ApiServerConfig struct {
	ApiPort              string        `yaml:"api_port"`
	AccessTokenDuration  time.Duration `yaml:"access_token_duration"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration"`
}

type DatabaseConfig struct {
	DbHost     string `yaml:"db_host"`
	DbPort     string `yaml:"db_port"`
	DbName     string `yaml:"db_name"`
	DbUsername string `yaml:"db_username"`
	DbPassword string `yaml:"db_password"`
}

type AWSConfig struct {
	Bucket          string `yaml:"bucket"`
	AccessKey       string `yaml:"access_key"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Region          string `yaml:"region"`
	BaseURL         string `yaml:"base_url"`
}

type OTPConfig struct {
	AccountSID string `yaml:"account_sid"`
	ApiKey     string `yaml:"api_key"`
	ApiSecret  string `yaml:"api_secret"`
	FromNumber string `yaml:"from_number"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Password string `yaml:"password"`
}

type PDFServiceConfig struct {
	Url     string        `yaml:"url"`
	Timeout time.Duration `yaml:"timeout"`
}

type BackgroundJobConfig struct {
	RentingBackoff time.Duration `yaml:"renting_backoff"`
}

type VNPayConfig struct {
	PayURL     string `yaml:"pay_url"`
	HashSecret string `yaml:"hash_secret"`
	Version    string `yaml:"version"`
	Command    string `yaml:"command"`
	TMNCode    string `yaml:"tmn_code"`
	Locale     string `yaml:"locale"`
	IpnURL     string `yaml:"ipn_url"`
	BankCode   string `yaml:"bank_code"`
}

type GlobalConfig struct {
	Database      *DatabaseConfig      `yaml:"database"`
	ApiServer     *ApiServerConfig     `yaml:"api_server"`
	AWS           *AWSConfig           `yaml:"aws"`
	OTP           *OTPConfig           `yaml:"otp"`
	PDFService    *PDFServiceConfig    `yaml:"pdf_service"`
	VNPay         *VNPayConfig         `yaml:"vn_pay"`
	BackgroundJob *BackgroundJobConfig `yaml:"background_job"`
	Redis         *RedisConfig         `yaml:"redis"`
}

type FEConfig struct {
	Path           string `yaml:"path"`
	AdminReturnURL string `yaml:"admin_return_url"`
	AdminBaseURL   string `yaml:"admin_base_url"`
}

func LoadConfig(path string) (*GlobalConfig, error) {
	cfg := &GlobalConfig{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	d := yaml.NewDecoder(file)
	if err = d.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func LoadFEConfig(path string) (*FEConfig, error) {
	cfg := &FEConfig{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	d := yaml.NewDecoder(file)
	if err = d.Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
