package admin

import (
	"github.com/lunny/tango"
	"github.com/tango-contrib/renders"
)

type Article struct {
	tango.Ctx
	renders.Renderer
}

func (a *Article) Get() {
	if err := a.Render("admin/article.tmpl"); err != nil {
		panic(err)
	}
}
