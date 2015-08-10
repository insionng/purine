package route

import (
	"github.com/fuxiaohei/purine/src/mapi"
	"github.com/fuxiaohei/purine/src/route/base"
	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
)

type Comment struct {
	base.BaseBinder

	tango.Ctx
	xsrf.Checker
}

func (c *Comment) Get() {
	c.WriteHeader(401)
}

func (c *Comment) Post() {
	form := new(mapi.CommentForm)
	if err := c.Bind(form); err != nil {
		c.ServeJson(mapi.Fail(err))
		return
	}
	form.For = c.Param("from")
	form.ForId = c.ParamInt64("id")
	// filter comment data
	if res := mapi.Call(mapi.Comment.Filter, form); !res.Status {
		c.ServeJson(res)
		return
	}
	// save comment
	res := mapi.Call(mapi.Comment.Save, form)
	c.ServeJson(res)
}
