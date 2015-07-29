package mapi

import (
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/utils"
	"path"
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
func (_ *RssApi) RSS(v interface{}) *Res {
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
	rss := &model.Rss{
		Title:       generalSetting.Title + " - " + generalSetting.Subtitle,
		Link:        generalSetting.BaseUrl,
		Description: generalSetting.Desc,
		Items:       make([]*model.Rss, len(articles)),
		PubDate:     time.Now(),
	}
	for i, a := range articles {
		r := &model.Rss{
			Title:       a.Title,
			Link:        path.Join(rss.Link, "article", a.Href()),
			Description: utils.Md2String(a.Body),
			PubDate:     time.Unix(a.CreateTime, 0),
			Items:       nil,
		}
		if i == 0 {
			rss.PubDate = r.PubDate
		}
		rss.Items[i] = r
	}
	return rss
}
