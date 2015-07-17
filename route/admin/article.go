package admin

import (
	"errors"
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/fuxiaohei/purine/utils"
	"github.com/lunny/tango"
)

type Article struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (a *Article) Get() {
	a.Title("Article")
	opt := &mapi.ArticleListOption{
		Page: a.FormInt64("page", 1),
	}
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

type Delete struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (d *Delete) Get() {
	id := d.FormInt64("id")
	// go back if no id
	if id == 0 {
		d.Redirect(d.Req().Referer())
		return
	}
	// get article
	article, err := model.GetArticleBy("id", id)
	if err != nil {
		panic(err)
	}
	if article.Id != id {
		d.Redirect(d.Req().Referer())
		return
	}
	d.Title("Delete - " + article.Title)
	d.Assign("Article", article)
	d.Render("delete.tmpl")
}

func (d *Delete) Post() {
	id := d.FormInt64("id")
	res := mapi.Call(mapi.DelArticle, id)
	if !res.Status {
		d.RenderError(errors.New(res.Error))
		return
	}
	d.Redirect("/admin/article")
}
