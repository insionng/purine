package mapi

import (
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/vars"
)

type ArticleListOption struct {
	OnlyStatus string
	Page       int
	Size       int
	Order      string
}

func prepareArticleListOption(opt *ArticleListOption) *ArticleListOption {
	if opt.Size == 0 {
		opt.Size = 10
	}
	if opt.Page == 0 {
		opt.Page = 1
	}
	if opt.Order == "" {
		opt.Order = "id DESC"
	}
	return opt
}

func ListArticle(v interface{}) *Res {
	opt, ok := v.(*ArticleListOption)
	if !ok {
		return Fail(paramTypeError(new(ArticleListOption)))
	}
	opt = prepareArticleListOption(opt)

	sess := vars.Db.NewSession()
	defer sess.Close()

	if opt.OnlyStatus != "" {
		sess.Where("status = ?", opt.OnlyStatus)
	} else {
		sess.Where("status != ?", model.ARTICLE_STATUS_DELETE)
	}
	sess.Limit(opt.Page, (opt.Page-1)*opt.Size).OrderBy(opt.Order)

	articles := []*model.Article{}
	if err := sess.Find(&articles); err != nil {
		log.Error("Db|ListArticle|%v|%s", opt, err.Error())
		return Fail(err)
	}

	return Success(map[string]interface{}{
		"articles": articles,
	})
}
