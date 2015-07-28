// Package model provides database and files operation
package model

import "github.com/fuxiaohei/purine/vars"

// main config
type Config struct {
	Version string
	Date    string
	Server  ConfigServer
}

// server config
type ConfigServer struct {
	Host string
	Port string
}

// new default config
func NewConfig() *Config {
	return &Config{
		Version: vars.VERSION,
		Date:    vars.VERSION_DATE,
		Server: ConfigServer{
			Host: "0.0.0.0",
			Port: "9999",
		},
	}
}
