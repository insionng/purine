package admin

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type Index struct {
	tango.Ctx
	renders.Renderer
}

func (i *Index) Get() {
	if err := i.Render("admin/index.tmpl"); err != nil {
		panic(err)
	}
}
