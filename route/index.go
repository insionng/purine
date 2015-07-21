package route

import (
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
	opt := &mapi.ArticleListOption{
		Page: idx.FormInt64("page", 1),
		Size: 4,
	}
	res := mapi.Call(mapi.ListArticle, opt)
	if !res.Status {
		panic(res.Error)
	}
	idx.Assign("Articles", res.Data["articles"].([]*model.Article))
	idx.Assign("Pager", res.Data["pager"].(*utils.Pager))
	idx.Render("index.tmpl")
}
