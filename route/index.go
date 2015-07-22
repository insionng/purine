package route

import (
	"errors"
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/fuxiaohei/purine/utils"
	"github.com/lunny/tango"
)

type Index struct {
	base.ThemeRender
	tango.Ctx
}

func (idx *Index) Get() {
	page := idx.ParamInt64(":page", 1)
	opt := &mapi.ArticleListOption{
		Page: page,
		Size: 2,
	}
	res := mapi.Call(mapi.ListArticle, opt)
	if !res.Status {
		idx.RenderError(500, errors.New(res.Error))
		return
	}
	idx.Assign("Articles", res.Data["articles"].([]*model.Article))
	idx.Assign("Pager", res.Data["pager"].(*utils.Pager))
	idx.Render("index.tmpl")
}
