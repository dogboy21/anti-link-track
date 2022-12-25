package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
)

type Config struct {
	configFile string `json:"-"`

	CleanedUrls            uint64 `json:"cleanedUrls"`
	AllowReferralMarketing bool   `json:"allowReferralMarketing"`
}

func (c *Config) Read() error {
	configBaseDir := ""
	switch runtime.GOOS {
	case "windows":
		configBaseDir = os.Getenv("APPDATA")
	case "darwin":
		configBaseDir = path.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		if os.Getenv("XDG_CONFIG_HOME") != "" {
			configBaseDir = os.Getenv("XDG_CONFIG_HOME")
		} else {
			configBaseDir = path.Join(os.Getenv("HOME"), ".config")
		}
	}

	if configBaseDir == "" {
		return fmt.Errorf("no config directory found for %s", runtime.GOOS)
	}

	c.configFile = path.Join(configBaseDir, "antilinktrack", "antilinktrack.json")

	os.MkdirAll(path.Dir(c.configFile), os.ModePerm)

	if _, err := os.Stat(c.configFile); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	configContent, err := os.ReadFile(c.configFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(configContent, c)
}

func (c *Config) Save() error {
	configContent, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(c.configFile, configContent, 0o755)
}
