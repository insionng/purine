package admin

import (
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/fuxiaohei/purine/utils"
	"github.com/lunny/tango"
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

type Write struct {
	base.AdminRender
	base.BaseBinder
	base.BaseAuther
	tango.Ctx
}

func (w *Write) Get() {
	w.Title("Write")
	w.Render("write.tmpl")
}

// ajax callback
func (w *Write) Post() {
	form := new(mapi.ArticleForm)
	if err := w.Bind(form); err != nil {
		w.ServeJson(mapi.Fail(err))
		return
	}
	form.AuthorId = w.AuthUser.Id

	res := mapi.Call(mapi.WriteArticle, form)
	w.ServeJson(res)
}
