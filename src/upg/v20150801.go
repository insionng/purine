package upg

import (
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/vars"
)

func init() {
	Script[20150801] = Update20150801
}

func Update20150801() error {
	// sync new models
	if err := vars.Db.Sync(new(model.Page), new(model.Comment)); err != nil {
		return err
	}
	return nil
}
