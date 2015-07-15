package admin

import (
	"github.com/fuxiaohei/purine/route/base"
)

type Index struct {
	base.AdminRender
}

func (i *Index) Get() {
	i.Title("Index")
	i.Render("index.tmpl")
}
