package core

import (
	"encoding/json"
	"os"
)

type Config struct {
	Token    string `json:"token"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Status   string `json:"status"`
}

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	AppContext.Config = &Config{}
	err = json.Unmarshal(data, AppContext.Config)
	if err != nil {
		return err
	}
	return nil
}
