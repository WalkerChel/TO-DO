package configs

import (
	"fmt"
	"path"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	HTTP `yaml:"http"`
	PG   `yaml:"db"`
}

type HTTP struct {
	Port string `yaml:"port" env:"HTTP_PORT"`
}

type PG struct {
	Username string `env-required:"true" yaml:"username"`
	Password string `env-required:"true" yaml:"password"`
	Host     string `env-required:"true" yaml:"host"`
	Port     string `env-required:"true" yaml:"port"`
	DBname   string `env-required:"true" yaml:"dbname"`
	SslMode  string `env-required:"true" yaml:"sslmode"`
}

func New(configPath string) (*Config, error) {
	cnf := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cnf)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %w", err)
	}

	err = cleanenv.UpdateEnv(cnf)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cnf, nil
}
