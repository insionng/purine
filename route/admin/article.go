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
	res := mapi.Call(mapi.Article.List, opt)
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
	if w.Form("type") == "page" {
		w.getPage()
		return
	}
	w.Title("Write")
	id := w.FormInt64("id")
	if id > 0 {
		article, err := model.GetArticleBy("id", id)
		if err != nil {
			w.RenderError(500, err)
			return
		}
		w.Title("Write - " + article.Title)
		w.Assign("Article", article)
	}
	w.Render("write.tmpl")
}

func (w *Write) getPage() {
	w.Title("Write Page")
	id := w.FormInt64("id")
	if id > 0 {
		page, err := model.GetPageBy("id", id)
		if err != nil {
			w.RenderError(500, err)
			return
		}
		w.Title("Write Page - " + page.Title)
		w.Assign("Page", page)
	}
	w.Render("write_page.tmpl")
}

// ajax callback
func (w *Write) Post() {
	if w.Form("type") == "page" {
		w.postPage()
		return
	}
	form := new(mapi.ArticleForm)
	if err := w.Bind(form); err != nil {
		w.ServeJson(mapi.Fail(err))
		return
	}
	form.AuthorId = w.AuthUser.Id

	res := mapi.Call(mapi.Article.Write, form)
	w.ServeJson(res)
}

func (w *Write) postPage() {
	form := new(mapi.PageForm)
	if err := w.Bind(form); err != nil {
		w.ServeJson(mapi.Fail(err))
		return
	}
	form.AuthorId = w.AuthUser.Id
	res := mapi.Call(mapi.Page.Write, form)
	w.ServeJson(res)
}

type Page struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (p *Page) Get() {
	opt := &mapi.PageListOption{
		Page: p.FormInt64("page", 1),
	}
	res := mapi.Call(mapi.Page.List, opt)
	if !res.Status {
		panic(res.Error)
	}
	p.Assign("Pages", res.Data["pages"].([]*model.Page))
	p.Assign("Pager", res.Data["pager"].(*utils.Pager))
	p.Title("Page")
	p.Render("page.tmpl")
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
	res := mapi.Call(mapi.Article.Remove, id)
	if !res.Status {
		d.RenderError(500, errors.New(res.Error))
		return
	}
	d.Redirect("/admin/article")
}
