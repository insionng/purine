package model

import "github.com/fuxiaohei/purine/vars"

type Config struct {
	Version string
	Date    string
	Server  ConfigServer
}

type ConfigServer struct {
	Host string
	Port string
}

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
