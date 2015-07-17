package mapi

import (
	"errors"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/utils"
	"strings"
)

const (
	ARTICLE_MORE_LINK = "<!--more-->"
)

var (
	ERR_ARTICLE_LINK = errors.New("article-link-exist")
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

type ArticleForm struct {
	Title    string `form:"title" binding:"Required"`
	Link     string `form:"link" binding:"Required;AlphaDashDot"`
	Body     string `form:"body" binding:"Required"`
	Format   string `form:"format"`
	Tag      string `form:"tag"`
	Draft    string `form:"draft"`
	Id       int64  `form:"id"`
	AuthorId int64
}

func WriteArticle(v interface{}) *Res {
	form, ok := v.(*ArticleForm)
	if !ok {
		return Fail(paramTypeError(new(ArticleForm)))
	}
	// check link
	a, err := model.GetArticleBy("link", form.Link)
	if err != nil {
		return Fail(err)
	}
	if a != nil {
		if form.Id == 0 || form.Id != a.Id {
			return Fail(ERR_ARTICLE_LINK)
		}
	}

	// create article object
	article := &model.Article{
		Id:            form.Id,
		AuthorId:      form.AuthorId,
		Title:         form.Title,
		Link:          form.Link,
		Body:          form.Body,
		TagString:     form.Tag,
		Hits:          1,
		Comments:      0,
		Status:        model.ARTICLE_STATUS_PUBLISH,
		CommentStatus: model.ARTICLE_COMMENT_OPEN,
	}
	if form.Draft != "" {
		article.Status = model.ARTICLE_STATUS_DRAFT
	}
	if strings.Contains(article.Body, ARTICLE_MORE_LINK) {
		article.Preview = strings.Split(article.Body, ARTICLE_MORE_LINK)[0]
	}

	// save article
	article, err = model.SaveArticle(article)
	if err != nil {
		return Fail(err)
	}

	return Success(map[string]interface{}{
		"article": article,
	})
}

func DelArticle(v interface{}) *Res {
	id, ok := v.(int64)
	if !ok {
		return Fail(paramTypeError(id))
	}
	if err := model.RemoveArticle(id); err != nil {
		return Fail(err)
	}
	return Success(nil)
}
