package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// ReadConfig reads configuration from given path.
func ReadConfig(configPath string) (*Config, error) {

	config := new(Config)

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	configReader := bytes.NewReader(configBytes)
	if err := json.NewDecoder(configReader).Decode(config); err != nil {
		return config, err
	}

	return config, nil
}

// Config is a struct representing whole application configuration.
type Config struct {
	Server   ServerConfig   `json:"server"`
	Session  SessionConfig  `json:"session"`
	Database DatabaseConfig `json:"database"`
	Statics  StaticsConfig  `json:"statics"`
}

// ServerConfig represents http server configuration.
type ServerConfig struct {
	Port int    `json:"port"`
	Host string `json:"host"`
}

// SessionConfig represents session configuration.
type SessionConfig struct {
	DB       int    `json:"db"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	IDLength int    `json:"id_len"`
}

// DatabaseConfig represents database configuration.
type DatabaseConfig struct {
	Host             string `json:"host"`
	Port             int    `json:"port"`
	DBName           string `json:"db_name"`
	UsersTableName   string `json:"users_table_name"`
	UsersTablePKName string `json:"users_table_pk_name"`
}

// StaticsConfig represents static files configuration.
type StaticsConfig struct {
	Path string `json:"path"`
}
