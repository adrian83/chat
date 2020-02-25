package config

import (
	"github.com/kelseyhightower/envconfig"
)

// ReadConfig reads configuration properties from Environment.
func ReadConfig(prefix string) (*Config, error) {
	var cfg Config
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Config is a struct representing whole application configuration.
type Config struct {
	ServerPort        int    `json:"serverPort" envconfig:"SERVER_PORT"`
	ServerHost        string `json:"serverHost" envconfig:"SERVER_HOST"`
	SessionDbName     int    `json:"sessionDbName" envconfig:"SESSION_DB_NAME"`
	SessionDbPassword string `json:"sessionDbPassword" envconfig:"SESSION_DB_PASSWORD"`
	SessionDbHost     string `json:"sessionDbHost" envconfig:"SESSION_DB_HOST"`
	SessionDbPort     int    `json:"sessionDbPort" envconfig:"SESSION_DB_PORT"`
	DatabaseHost      string `json:"databaseHost" envconfig:"DATABASE_HOST"`
	DatabasePort      int    `json:"databasePort" envconfig:"DATABASE_PORT"`
	DatabaseName      string `json:"databaseName" envconfig:"DATABASE_NAME"`
	StaticsPath       string `json:"staticsPath" envconfig:"STATICS_PATH"`
}
