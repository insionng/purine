package route

import (
	"github.com/fuxiaohei/purine/src/mapi"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/route/base"
	"github.com/lunny/tango"
	"github.com/tango-contrib/xsrf"
)

var (
	blogParamRule = "id/link"
	pageParamRule = "link"
)

type Article struct {
	base.ThemeRender
	tango.Ctx
	xsrf.Checker
}

func (a *Article) Get() {
	param, err := mapi.Article.ParseRoute(blogParamRule, a.Param("*article"))
	if err != nil {
		a.RenderError(500, err)
		return
	}

	// get article
	res := mapi.Call(mapi.Article.Get, param)
	if res.Status {
		article := res.Data["article"].(*model.Article)
		if article.Status != model.ARTICLE_STATUS_PUBLISH {
			a.RenderError(404, nil)
			return
		}
		a.Assign("Article", article)
		a.Assign("Title", article.Title+" - "+a.GetGeneralByKey("title"))

		// get comments
		opt := &mapi.CommentListOption{
			ArticleId: article.Id,
			Page:      1,
			Size:      999,
			Status:    model.COMMENT_STATUS_APPROVED,
		}
		res = mapi.Call(mapi.Comment.List, opt)
		if res.Status {
			comments := res.Data["comments"].([]*model.Comment)
			a.Assign("Comments", comments)
		}
		a.Assign("XsrfHtml", a.XsrfFormHtml())
		a.Render("article.tmpl")
		return
	}

	// not found
	a.RenderError(404, nil)
}

type Page struct {
	base.ThemeRender
	tango.Ctx
}

func (p *Page) Get() {
	param, err := mapi.Page.ParseRoute(pageParamRule, p.Param("*page"))
	if err != nil {
		p.RenderError(500, err)
		return
	}
	res := mapi.Call(mapi.Page.Get, param)
	if res.Status {
		page := res.Data["page"].(*model.Page)
		if page.Status != model.PAGE_STATUS_PUBLISH {
			p.RenderError(404, nil)
			return
		}
		p.Assign("Page", page)
		p.Title(page.Title + " - " + p.GetGeneralByKey("title"))
		p.Render("page.tmpl")
		return
	}
	p.RenderError(404, nil)
}
