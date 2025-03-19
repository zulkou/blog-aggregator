package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
    DBURL               string  `json:"db_url"`
    CurrentUserName     string  `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("Home dir not found: %w\n", err)
    }
    
    return fmt.Sprintf("%v/%s", homeDir, configFileName), nil
}

func Read() (Config, error) {
    path, err := getConfigFilePath()
    if err != nil {
        return Config{}, fmt.Errorf("Error in retrieving filepath: %w\n", err)
    }

    file, err := os.Open(path)
    if err != nil {
        return Config{}, fmt.Errorf("Error opening file: %w\n", err)
    }
    defer file.Close()

    data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, fmt.Errorf("Error on reading file: %w\n", err)
    }

    var res Config
    err = json.Unmarshal(data, &res)
    if err != nil {
        return Config{}, fmt.Errorf("Error on parsing file: %w\n", err)
    }

    return res, nil
}

func write(cfg Config) error {
    path, err := getConfigFilePath()
    if err != nil {
        return fmt.Errorf("Error in retrieving filepath: %w", err)
    }

    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("Error opening file: %w", err)
    }
    defer file.Close()

    jsonData, err := json.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("Error marshalling data into json: %w", err)
    }

    return os.WriteFile(path, jsonData, 0644)
}

func (c *Config) SetUser(user string) error {
    c.CurrentUserName = user
    return write(*c)
}
