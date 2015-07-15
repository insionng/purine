package admin

import (
	"github.com/fuxiaohei/purine/route/base"
)

type Article struct {
	base.AdminRender
	base.BaseAuther
}

func (a *Article) Get() {
	a.Title("Article")
	a.Render("article.tmpl")
}
