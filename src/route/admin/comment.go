package admin

import (
	"github.com/fuxiaohei/purine/src/route/base"
	"github.com/lunny/tango"
)

type Comment struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (c *Comment) Get() {
	c.Title("Comments")
	c.Render("comment.tmpl")
}
