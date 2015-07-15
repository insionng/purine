package model

import (
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/vars"
	"net/url"
	"time"
)

const (
	ARTICLE_STATUS_PUBLISH = "publish"
	ARTICLE_STATUS_DRAFT   = "draft"
	ARTICLE_STATUS_DELETE  = "deleted"

	ARTICLE_COMMENT_OPEN  = "open"
	ARTICLE_COMMENT_CLOSE = "closed"
	ARTICLE_COMMENT_WAIT  = "wait" // waiting 30 days to close
)

type Article struct {
	Id       int64
	AuthorId int64

	Title      string
	CreateTime int64  `xorm:"created"`
	UpdateTime int64  `xorm:"updated"`
	Link       string `xorm:"unique"`
	Preview    string
	Body       string
	Topic      string
	TagString  string

	Hits          int64
	Comments      int64
	Status        string
	CommentStatus string

	tagData []*Tag `xorm:"-"`
}

func (a *Article) Date(layout string) string {
	return time.Unix(a.CreateTime, 0).Format(layout)
}

func (a *Article) Href() string {
	return url.QueryEscape(a.Title)
}

func (a *Article) IsDraft() bool {
	return a.Status == ARTICLE_STATUS_DRAFT
}

func ListGeneralArticle(page, size int64, order string) ([]*Article, error) {
	articles := make([]*Article, 0)
	if err := vars.Db.Where("status != ?", ARTICLE_STATUS_DELETE).
		Limit(int(page), int((page-1)*size)).OrderBy(order).Find(&articles); err != nil {
		log.Error("Db|ListGeneralArticle|%d,%d|%s|%s", page, size, order, err.Error())
		return nil, err
	}
	return articles, nil
}

func CountGeneralArticle() (int64, error) {
	return vars.Db.Where("status != ?", ARTICLE_STATUS_DELETE).Count(new(Article))
}

func ListStatusArticle(status string, page, size int64, order string) ([]*Article, error) {
	articles := make([]*Article, 0)
	if err := vars.Db.Where("status = ?", status).
		Limit(int(page), int((page-1)*size)).OrderBy(order).Find(&articles); err != nil {
		log.Error("Db|ListStatusArticle|%s|%d,%d|%s|%s", status, page, size, order, err.Error())
		return nil, err
	}
	return articles, nil
}

func CountStatusArticle(status string) (int64, error) {
	return vars.Db.Where("status = ?", status).Count(new(Article))
}

type Tag struct {
	Id        int64
	ArticleId int64
	Tag       string
}
