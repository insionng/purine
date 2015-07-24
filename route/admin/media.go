package admin

import "github.com/fuxiaohei/purine/route/base"

type Media struct {
	base.AdminRender
	base.BaseAuther
	base.BaseBinder
}

func (m *Media) Get() {
	m.Title("Media")
	m.Render("media.tmpl")
}
