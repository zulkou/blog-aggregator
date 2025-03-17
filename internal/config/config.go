package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
    DBURL              string  `json:"db_url"`
    CurrentUserName   string  `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", fmt.Errorf("home dir not found: %w", err)
    }
    
    return fmt.Sprintf("%v/%s", homeDir, configFileName), nil
}

func Read() (Config, error) {
    path, err := getConfigFilePath()
    if err != nil {
        return Config{}, fmt.Errorf("path error: %w", err)
    }

    file, err := os.Open(path)
    if err != nil {
        return Config{}, fmt.Errorf("error opening file: %w", err)
    }
    defer file.Close()

    data, err := os.ReadFile(path)
    if err != nil {
        return Config{}, fmt.Errorf("error on reading file: %w", err)
    }

    var res Config
    err = json.Unmarshal(data, &res)

    return res, nil
}

func write(cfg Config) error {
    path, err := getConfigFilePath()
    if err != nil {
        return fmt.Errorf("path error: %w", err)
    }

    file, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("error opening file: %w", err)
    }
    defer file.Close()

    jsonData, err := json.Marshal(cfg)
    if err != nil {
        return fmt.Errorf("error json marshal: %w", err)
    }

    return os.WriteFile(path, jsonData, 0644)
}

func (c *Config) SetUser(user string) error {
    c.CurrentUserName = user
    return write(*c)
}
