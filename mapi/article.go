package mapi

import (
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/utils"
)

type ArticleListOption struct {
	Status string
	Page   int64
	Size   int64
	Order  string
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

	if opt.Status != "" {
		articles, err := model.ListStatusArticle(opt.Status, opt.Page, opt.Size, opt.Order)
		if err != nil {
			return Fail(err)
		}
		count, err := model.CountStatusArticle(opt.Status)
		if err != nil {
			return Fail(err)
		}
		return Success(map[string]interface{}{
			"articles": articles,
			"pager":    utils.CreatePager(opt.Page, opt.Size, count),
		})
	}

	articles, err := model.ListGeneralArticle(opt.Page, opt.Size, opt.Order)
	if err != nil {
		return Fail(err)
	}
	count, err := model.CountGeneralArticle()
	if err != nil {
		return Fail(err)
	}
	return Success(map[string]interface{}{
		"articles": articles,
		"pager":    utils.CreatePager(opt.Page, opt.Size, count),
	})
}
