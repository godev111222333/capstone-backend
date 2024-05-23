package misc

import (
	"gopkg.in/yaml.v3"
	"os"
)

type DatabaseConfig struct {
	DbHost     string `yaml:"db_host"`
	DbPort     string `yaml:"db_port"`
	DbName     string `yaml:"db_name"`
	DbUsername string `yaml:"db_username"`
	DbPassword string `yaml:"db_password"`
}

type GlobalConfig struct {
	Database *DatabaseConfig
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
