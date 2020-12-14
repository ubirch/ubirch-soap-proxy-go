package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Uuid                string `json:"uuid"`
	Auth                string `json:"auth"`
	VerificationBaseURL string `json:"verificationBaseURL"`
	UbirchClientURL     string `json:"ubirchClientURL"`
}

func (c *Config) Load(configDir string, filename string) error {
	if os.Getenv("UBIRCH_UBIRCH_CLIENT_URL") != "" {
		return c.loadEnv()
	} else {
		return c.loadFile(filepath.Join(configDir, filename))
	}
}

// loadEnv reads the configuration from environment variables
func (c *Config) loadEnv() error {
	return envconfig.Process("ubirch", c)
}

// LoadFile reads the configuration from a json file
func (c *Config) loadFile(filename string) error {
	contextBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(contextBytes, c)
}
