package database

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrConfigInvalidSSLMode = errors.New("invalid SSLMode; valid options are [require, disable]")

type SSLMode string

const (
	SSLModeRequire SSLMode = "require"
	SSLModeDisable SSLMode = "disable"
)

var (
	SSLModes = map[SSLMode]struct{}{
		SSLModeRequire: {},
		SSLModeDisable: {},
	}
)

// Config describes necessary information
// to access a Postgres instance
type Config struct {
	Host     string  `mapstructure:"host" yaml:"host"`
	Port     int64   `mapstructure:"port" yaml:"port"`
	Database string  `mapstructure:"database" yaml:"database"`
	Username string  `mapstructure:"username" yaml:"username"`
	Password string  `mapstructure:"password" yaml:"password"`
	SSLMode  SSLMode `mapstructure:"sslmode" yaml:"sslmode"`
}

func Defaults() *Config {
	return &Config{
		Host:     "0.0.0.0",
		Port:     5432,
		Database: "postgres",
		Username: "postgres",
		Password: "postgres",
		SSLMode:  SSLModeDisable,
	}
}

// marshalJSON unmarshals a JSON string into a SSLMode and
// checks if it's a valid supported option [require, disable]
func (s *SSLMode) UnmarshalJSON(data []byte) error {
	var mode string
	if err := json.Unmarshal(data, &mode); err != nil {
		return fmt.Errorf("[database] failed unmarshalling database config: %w", err)
	}

	ssm := SSLMode(mode)
	if _, valid := SSLModes[ssm]; !valid {
		return fmt.Errorf("[database] failed validating database config: %w", ErrConfigInvalidSSLMode)
	}

	*s = ssm
	return nil
}
