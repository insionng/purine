package admin

import (
	"github.com/fuxiaohei/purine/src/route/base"
)

type Index struct {
	base.AdminRender
	base.BaseAuther
}

func (i *Index) Get() {
	i.Title("Index")
	i.Render("index.tmpl")
}
