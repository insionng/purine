package admin

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type Login struct {
	renders.Renderer
	tango.Ctx
}

func (l *Login) Get() {
	if err := l.Render("admin/login.tmpl"); err != nil {
		panic(err)
	}
}
