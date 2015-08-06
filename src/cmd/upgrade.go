package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/src/log"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/upg"
	"github.com/fuxiaohei/purine/src/vars"
	"sort"
	"strconv"
	"time"
)

var upgradeCmd cli.Command = cli.Command{
	Name:  "upgrade",
	Usage: "upgrade blog framework",
	Action: func(ctx *cli.Context) {
		cfg, err := loadConfig()
		if err != nil {
			log.Error("Upgrade | %-8s | %s", "Config", err.Error())
			return
		}
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

	opt := &PrepareOption{true, true, true}
	pre, err := Prepare(opt)
	if err != nil {
		log.Fatal("Upgrade | %-8s | %s", "Prepare", err.Error())
		return
	}

	oldVersion, _ := strconv.Atoi(pre.Config.Date)
	scriptIndex := []int{}
	for vr, _ := range upg.Script {
		if vr > oldVersion {
			scriptIndex = append(scriptIndex, vr)
		}
	}
	sort.Sort(sort.IntSlice(scriptIndex))

	for _, cv := range scriptIndex {
		log.Debug("Upgrade | %-8s | %d ", "Process", cv)
		if err := upg.Script[cv](); err != nil {
			log.Error("Upgrade | %-8s | %s", "Process", err.Error())
			return
		}
	}

	pre.Config.Version = vars.VERSION
	pre.Config.Date = vars.VERSION_DATE
	if err := model.SyncConfig(pre.Config); err != nil {
		log.Error("Upgrade | %-8s | SyncFail", "Config")
		return
	}

	log.Info("Upgrade | %-8s | Sync | %s", "Config", vars.CONFIG_FILE)
	log.Info("Upgrade | %-8s | %.1fms", "Done", time.Since(t).Seconds()*1000)
}
