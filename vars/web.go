package vars

import (
	"github.com/fuxiaohei/purine/log"
	"github.com/lunny/tango"
)

var Server *tango.Tango

func init() {
	Server = tango.NewWithLog(log.Get().ToTangoLogger(), []tango.Handler{
		tango.Logging(),
		tango.Recovery(true),
		tango.Return(),
		tango.Param(),
		tango.Contexts(),
	}...)
}
