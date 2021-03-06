package route

import (
	"errors"
	"github.com/fuxiaohei/purine/src/mapi"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/route/base"
	"github.com/fuxiaohei/purine/src/utils"
	"github.com/lunny/tango"
	"strconv"
)

type Index struct {
	base.ThemeRender
	tango.Ctx
}

func (idx *Index) Get() {
	page := idx.ParamInt64(":page", 1)
	if page > 1 {
		idx.Assign("Title", idx.GetGeneralByKey("title")+" - Page "+strconv.FormatInt(page, 10))
	}
	opt := &mapi.ArticleListOption{
		Page:   page,
		Size:   4,
		Status: model.ARTICLE_STATUS_PUBLISH,
	}
	res := mapi.Call(mapi.Article.List, opt)
	if !res.Status {
		idx.RenderError(500, errors.New(res.Error))
		return
	}
	pager := res.Data["pager"].(*utils.Pager)
	if pager.Current > pager.Pages {
		idx.RenderError(404, nil)
		return
	}
	idx.Assign("Articles", res.Data["articles"].([]*model.Article))
	idx.Assign("Pager", pager)
	idx.Render("index.tmpl")
}

type Archive struct {
	base.ThemeRender
	tango.Ctx
}

func (a *Archive) Get() {
	a.Title("Archive")
	res := mapi.Call(mapi.Article.ListArchive, nil)
	if !res.Status {
		a.RenderError(500, errors.New(res.Error))
		return
	}
	a.Assign("Articles", res.Data["articles"].([]*model.Article))
	a.Render("archive.tmpl")
}
