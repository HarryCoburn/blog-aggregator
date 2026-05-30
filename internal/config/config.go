package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("Cannot get user's home directory string.")
	}
	path := filepath.Join(home_dir, configFileName)
	return path, nil
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	data, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("Could not open file.")
	}
	defer data.Close()
	var config Config
	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("Could not decode JSON data: ")
	}
	return config, nil
}

func (c Config) SetUser(user string) {
	c.Current_user_name = user
	write(c)
}

func write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	data, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Could not open file.")
	}
	defer data.Close()
	if err := json.NewEncoder(data).Encode(cfg); err != nil {
		return fmt.Errorf("Could not write JSON data to config file.")
	}
	return nil
}
