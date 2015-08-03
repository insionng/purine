package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/upg"
	"github.com/fuxiaohei/purine/vars"
	"sort"
	"strconv"
	"time"
)

var upgradeCmd cli.Command = cli.Command{
	Name:  "upgrade",
	Usage: "upgrade blog framework",
	Action: func(ctx *cli.Context) {
		cfg := UpgradeConfig(ctx)
		if cfg == nil {
			log.Error("Upgrade | %-8s | ReadFail", "Config")
			return
		}
		if IsNeedUpgrade(cfg) {
			UpgradeAction(cfg)
			return
		}

		log.Error("Upgrade | %-8s | Same Version", "Config")

	},
}

// check if need upgrade
func IsNeedUpgrade(cfg *model.Config) bool {
	if cfg.Version != vars.VERSION {
		return true
	}
	return vars.VERSION_DATE != cfg.Date
}

// upgrade action
func UpgradeAction(cfg *model.Config) {
	t := time.Now()
	log.Debug("Upgrade | %-8s | %s -> %s", "Upgrade", cfg.Version, vars.VERSION)

	oldVersion, _ := strconv.Atoi(cfg.Date)
	scriptIndex := []int{}
	for vr, _ := range upg.Script {
		if vr > oldVersion {
			scriptIndex = append(scriptIndex, vr)
		}
	}
	sort.Sort(sort.IntSlice(scriptIndex))
	if err := prepareUpgrade(); err != nil {
		log.Error("Upgrade | %-8s | Fail", "Prepare")
		return
	}

	for _, cv := range scriptIndex {
		log.Debug("Upgrade | %-8s | %d ", "Process", cv)
		if err := upg.Script[cv](); err != nil {
			log.Error("Upgrade | %-8s | %s", "Process", err.Error())
			return
		}
	}

	cfg.Version = vars.VERSION
	cfg.Date = vars.VERSION_DATE
	if err := model.SyncConfig(cfg); err != nil {
		log.Error("Upgrade | %-8s | SyncFail", "Config")
		return
	}

	log.Info("Upgrade | %-8s | Sync | %s", "Config", vars.CONFIG_FILE)
	log.Info("Upgrade | %-8s | %.1fms", "Done", time.Since(t).Seconds()*1000)
}

// prepare upgrade processes
func prepareUpgrade() error {
	return loadDb()
}

// load config to server
func UpgradeConfig(ctx *cli.Context) *model.Config {
	cfg, err := loadConfig()
	if err != nil {
		log.Error("Upgrade | %-8s | %s", "Config", err.Error())
		return nil
	}
	return cfg
}
