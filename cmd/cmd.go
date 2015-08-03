// Package cmd contains all commands
package cmd

import (
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/vars"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	vars.Cli.Commands = []cli.Command{
		installCmd,
		versionCmd,
		servCmd,
		packCmd,
		upgradeCmd,
	}
}

func loadConfig() (*model.Config, error) {
	cfg := new(model.Config)
	if _, err := toml.DecodeFile(vars.CONFIG_FILE, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadDb() error {
	engine, err := xorm.NewEngine("sqlite3", vars.DATA_FILE)
	if err != nil {
		return err
	}
	engine.SetLogger(nil) // close logger
	vars.Db = engine
	return nil
}
