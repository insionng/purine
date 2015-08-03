// Package model provides database and files operation
package model

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/fuxiaohei/purine/vars"
	"io/ioutil"
	"os"
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

func SyncConfig(cfg *Config) error {
	// encode config
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(cfg); err != nil {
		return err
	}
	// write config.toml
	if err := ioutil.WriteFile(vars.CONFIG_FILE, buf.Bytes(), os.ModePerm); err != nil {
		return err
	}
	return nil
}
