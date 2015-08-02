package admin

import (
	"github.com/fuxiaohei/purine/route/base"
	"github.com/lunny/tango"
)

type Page struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (p *Page) Get() {
	p.Title("Page")
	p.Render("page.tmpl")
}
