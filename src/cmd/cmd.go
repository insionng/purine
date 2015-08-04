// Package cmd contains all commands
package cmd

import (
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/vars"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
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

// prepare option
type PrepareOption struct {
	LoadConfig bool
	LoadDb     bool
	LoadI18n   bool
}

// prepare loaded data
type PreparedData struct {
	Config *model.Config
}

// merge default prepare options
func mergePrepareOption(opt *PrepareOption) *PrepareOption {
	opt.LoadConfig = true
	return opt
}

func Prepare(opt *PrepareOption) (*PreparedData, error) {
	opt = mergePrepareOption(opt)
	var (
		err  error
		data = new(PreparedData)
	)
	if opt.LoadConfig {
		if data.Config, err = loadConfig(); err != nil {
			return nil, err
		}
	}
	if opt.LoadDb {
		if err = loadDb(); err != nil {
			return nil, err
		}
	}
	if opt.LoadI18n {
		if err = loadI18n(); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// load config
func loadConfig() (*model.Config, error) {
	cfg := new(model.Config)
	if _, err := toml.DecodeFile(vars.CONFIG_FILE, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// load database
func loadDb() error {
	engine, err := xorm.NewEngine("sqlite3", vars.DATA_FILE)
	if err != nil {
		return err
	}
	engine.SetLogger(nil) // close logger
	vars.Db = engine
	return nil
}

func loadI18n() error {
	filepath.Walk(vars.I18N_DIR, func(path string, info os.FileInfo, err error) error {
		return nil
	})
	return nil
}
