package config

import (
	"encoding/json"
	"os"
)

const fileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return Config{}, err
	}
	var newConfig Config
	if err = json.Unmarshal(data, &newConfig); err != nil {
		return Config{}, err
	}
	return newConfig, nil
}

func (cfg *Config) SetUser(name string) error {
	cfg.CurrentUserName = name
	return write(*cfg)
}

func write(cfg Config) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(cfg); err != nil {
		return err
	}
	return nil
}


