package mapi

import (
	"github.com/fuxiaohei/purine/model"
	"time"
)

var (
	Rss = new(RssApi) // rss group api
)

// rss group api struct
type RssApi struct{}

// rss option
type RssOption struct {
	Setting  *SettingGeneral
	Articles []*model.Article
}

// create rss data
//
//  in  : *RssOption
//  out : {
//          "rss":*Rss
//        }
//
func (_ *RssApi) ListRSS(v interface{}) *Res {
	opt, ok := v.(*RssOption)
	if !ok {
		return Fail(paramTypeError(opt))
	}
	rss := articles2Rss(opt.Setting, opt.Articles)
	return Success(map[string]interface{}{
		"rss": rss,
	})
}

// change article to rss items
func articles2Rss(generalSetting *SettingGeneral, articles []*model.Article) *model.Rss {
	rss := model.Articles2Rss(articles, generalSetting.BaseUrl)
	rss.Title = generalSetting.Title + " - " + generalSetting.Subtitle
	rss.Description = generalSetting.Desc
	rss.PubDate = time.Now()
	if len(rss.Items) > 0 {
		rss.PubDate = rss.Items[0].PubDate
	}
	return rss
}

type SiteMapOption struct {
	Setting  *SettingGeneral
	Articles []*model.Article
}

func (_ *RssApi) ListSitemap(v interface{}) *Res {
	opt, ok := v.(*SiteMapOption)
	if !ok {
		return Fail(paramTypeError(opt))
	}
	sitemap := sitemap2Rss(opt.Setting, opt.Articles)
	return Success(map[string]interface{}{
		"sitemap": sitemap,
	})
}

func sitemap2Rss(generalSetting *SettingGeneral, articles []*model.Article) *model.SiteMapGroup {
	return model.Articles2SiteMap(articles, generalSetting.BaseUrl)
}
