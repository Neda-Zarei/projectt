package config

import (
	"encoding/json"
	"os"

	"github.com/caarlos0/env/v11"
)

func ReadEnv() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func ReadJson(file string) (cfg Config, err error) {
	var data []byte
	if data, err = os.ReadFile(file); err != nil {
		return Config{}, err
	}
	if err = json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, err
}
