package config

import (
	"os"

	"github.com/vandi37/vanerrors"
	"gopkg.in/yaml.v3"
)

// Errors
const (
	ErrorToOpenConfig = "error to open config"
	ErrorDecodingData = "error decoding data"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// Config
type Config struct {
	Port  int      `yaml:"port"`
	Token string   `yaml:"token"`
	DB    DBConfig `yaml:"db"`
}

// Gets config
func Get(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorToOpenConfig, err, vanerrors.EmptyHandler)
	}
	defer file.Close()

	cfg := new(Config)

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, vanerrors.NewWrap(ErrorDecodingData, err, vanerrors.EmptyHandler)
	}

	return cfg, nil
}
