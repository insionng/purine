package upg

import (
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/vars"
)

func init() {
	Script[20150801] = Update20150801
}

func Update20150801() error {
	if err := vars.Db.Sync(new(model.Page)); err != nil {
		return err
	}
	return nil
}
