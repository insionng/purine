package vars

import (
	"github.com/fuxiaohei/purine/src/log"
	"github.com/lunny/tango"
)

var (
	Server          *tango.Tango
	StaticDirectory map[string]string = map[string]string{
		"/static": "static",
	}
)

func init() {
	Server = tango.NewWithLog(log.Get().ToTangoLogger(), []tango.Handler{
		tango.Recovery(true),
		tango.Return(),
		tango.Param(),
		tango.Contexts(),
	}...)
}
