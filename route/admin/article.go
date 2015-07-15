package admin

import (
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/fuxiaohei/purine/utils"
)

type Article struct {
	base.AdminRender
	base.BaseAuther
}

func (a *Article) Get() {
	a.Title("Article")
	opt := &mapi.ArticleListOption{}
	res := mapi.Call(mapi.ListArticle, opt)
	if !res.Status {
		panic(res.Error)
	}
	a.Assign("Articles", res.Data["articles"].([]*model.Article))
	a.Assign("Pager", res.Data["pager"].(*utils.Pager))
	a.Render("article.tmpl")
}
