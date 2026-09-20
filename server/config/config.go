package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// vulnerable to cve-2020-10136 could potentially be used as relay!
type Config struct {
	Servers    []string `json:"servers"`
	ServerAddr string   `json:"server_addr"`
	RelayAddr  string   `json:"relay_addr"`
	WebListen  string   `json:"web_listen"`
	AdminKey   string   `json:"admin_key"`
}

func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file error: %w", err)
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	return &cfg, nil
}
