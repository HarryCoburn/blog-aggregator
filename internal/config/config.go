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

func Read() (Config, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return Config{}, errors.New("Cannot get user's home directory string.")
	}
	path := filepath.Join(home_dir, configFileName)
	data, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("Could not read file.")
	}
	defer data.Close()
	var config Config
	decoder := json.NewDecoder(data)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("Could not decode JSON data: ")
	}
	return config, nil
}
