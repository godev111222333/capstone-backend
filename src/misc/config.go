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
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type GlobalConfig struct {
	Database  *DatabaseConfig  `yaml:"database"`
	ApiServer *ApiServerConfig `yaml:"api_server"`
	AWS       *AWSConfig       `yaml:"aws"`
	OTP       *OTPConfig       `yaml:"otp"`
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
