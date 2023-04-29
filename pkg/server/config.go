package server

import (
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	DBHost     string `yaml:"db_host"`
	DBPort     int    `yaml:"db_port"`
	DBName     string `yaml:"db_name"`
	DBSchema   string `yaml:"db_schema"`
	DBUser     string `yaml:"db_user"`
	DBPassword string `yaml:"db_password"`

	TgToken string `yaml:"tg_token"`
}

func LoadConfig(pathToFile string) (*Config, error) {
	fileBody, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, xerrors.Errorf("unable to read file %s: %w", pathToFile, err)
	}

	var result Config
	err = yaml.Unmarshal(fileBody, &result)
	if err != nil {
		return nil, xerrors.Errorf("unable to unmarshal config: %w", err)
	}
	return &result, err
}
