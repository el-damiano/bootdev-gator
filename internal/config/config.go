package config

import (
	"encoding/json"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DatabaseUrl string `json:"db_url"`
	UserCurrent string `json:"current_user_name"`
}

func getConfigPath() (string, error) {
	configPath, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath += "/" + configFileName
	return configPath, nil
}

func Read() (Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, err
	}

	var config = Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (config *Config) SetUser(username string) error {
	config.UserCurrent = username

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
