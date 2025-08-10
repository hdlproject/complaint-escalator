package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Interval time.Duration `yaml:"interval"`
	Backoff  time.Duration `yaml:"backoff"`
	Template string        `yaml:"template"`
	Channels []string      `yaml:"channels"`
	// Azure Communication Services configuration
	ACS struct {
		ConnectionString string `yaml:"connection_string"`
		Domain           string `yaml:"domain"`
		FromEmail        string `yaml:"from_email"`
	} `yaml:"acs"`
	// Email configuration
	Email struct {
		To      []string `yaml:"to"`
		CC      []string `yaml:"cc,omitempty"`
		BCC     []string `yaml:"bcc,omitempty"`
		ReplyTo string   `yaml:"reply_to,omitempty"`
	} `yaml:"email"`
}

func LoadConfig(path string) (Config, error) {
	var cfg Config
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
