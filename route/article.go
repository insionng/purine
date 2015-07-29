package route

import (
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/lunny/tango"
)

var (
	blogParamRule = "id/link"
)

type Article struct {
	base.ThemeRender
	tango.Ctx
}

func (a *Article) Get() {
	param, err := mapi.Article.ParseRoute(blogParamRule, a.Param("*article"))
	if err != nil {
		a.RenderError(500, err)
		return
	}
	res := mapi.Call(mapi.Article.Get, param)
	if res.Status {
		article := res.Data["article"].(*model.Article)
		a.Assign("Article", article)
		a.Assign("Title", article.Title+" - "+a.GetGeneralByKey("title"))
		a.Render("article.tmpl")
		return
	}
	a.RenderError(404, nil)
}
