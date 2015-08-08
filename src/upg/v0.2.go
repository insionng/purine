package upg

import (
	"github.com/fuxiaohei/purine/src/log"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/vars"
	"reflect"
)

func init() {
	Script[20150807] = Update20150807
}

func Update20150807() error {
	key := "20150807"
	// sync new models
	if err := vars.Db.Sync(new(model.Page), new(model.Comment)); err != nil {
		return err
	}
	log.Debug("%s| %-8s | %s,%s", key, "SyncDb",
		reflect.TypeOf(new(model.Page)).String(),
		reflect.TypeOf(new(model.Comment)).String(),
	)
	return nil
}
