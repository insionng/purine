package route

import (
	"errors"
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/lunny/tango"
)

type Rss struct {
	base.BaseSetting
	tango.Ctx
}

func (r *Rss) Get() {
	res := mapi.Call(mapi.Article.ListArchive, nil)
	if !res.Status {
		panic(errors.New(res.Error))
		return
	}
	opt := &mapi.RssOption{
		Setting:  r.GetGeneral(),
		Articles: res.Data["articles"].([]*model.Article),
	}

	res = mapi.Call(mapi.Rss.ListRSS, opt)
	rss := res.Data["rss"].(*model.Rss)
	r.WriteHeader(200)
	r.Header().Set("Content-Type", "application/rss+xml;charset=UTF-8")
	r.Write([]byte(rss.String()))
}

type Sitemap struct {
	base.BaseSetting
	tango.Ctx
}

func (s *Sitemap) Get() {
	res := mapi.Call(mapi.Article.ListArchive, nil)
	if !res.Status {
		panic(errors.New(res.Error))
		return
	}
	opt := &mapi.SiteMapOption{
		Setting:  s.GetGeneral(),
		Articles: res.Data["articles"].([]*model.Article),
	}

	res = mapi.Call(mapi.Rss.ListSitemap, opt)
	rss := res.Data["sitemap"].(*model.SiteMapGroup)
	s.WriteHeader(200)
	s.Header().Set("Content-Type", "text/xml; charset=utf-8")
	s.Write([]byte(rss.String()))
}
