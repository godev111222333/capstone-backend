package misc

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ApiServerConfig struct {
	ApiPort string `yaml:"api_port"`
}

type DatabaseConfig struct {
	DbHost     string `yaml:"db_host"`
	DbPort     string `yaml:"db_port"`
	DbName     string `yaml:"db_name"`
	DbUsername string `yaml:"db_username"`
	DbPassword string `yaml:"db_password"`
}

type AWSConfig struct {
	AccessKey       string `yaml:"access_key"`
	SecretAccessKey string `yaml:"secret_access_key"`
}
type GlobalConfig struct {
	Database  *DatabaseConfig  `yaml:"database"`
	ApiServer *ApiServerConfig `yaml:"api_server"`
	AWS       *AWSConfig       `yaml:"aws"`
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
