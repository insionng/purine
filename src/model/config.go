// Package model provides database and files operation
package model

import (
	"github.com/fuxiaohei/purine/src/vars"
	"gopkg.in/ini.v1"
)

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

// write config to ini file
func WriteConfig(cfg *Config, file string) error {
	f := ini.Empty()
	if err := f.ReflectFrom(cfg); err != nil {
		return err
	}
	if err := f.SaveToIndent(file, "  "); err != nil {
		return err
	}
	f = nil
	return nil
}

// read config from file
func ReadConfig(file string) (*Config, error) {
	f, err := ini.Load(file)
	if err != nil {
		return nil, err
	}
	cfg := new(Config)
	if err = f.MapTo(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
