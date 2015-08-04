package model

import (
	"fmt"
	"github.com/fuxiaohei/purine/src/log"
	"github.com/fuxiaohei/purine/src/vars"
	"net/url"
	"strings"
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

// article struct
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

	tagData  []*Tag `xorm:"-"`
	userData *User  `xorm:"-"`
}

// format article create time
func (a *Article) Date(layout string) string {
	return time.Unix(a.CreateTime, 0).Format(layout)
}

// article visitable link
func (a *Article) Href() string {
	l := a.Link
	if l == "" {
		l = url.QueryEscape(a.Title)
	}
	return fmt.Sprintf("%d/%s.html", a.Id, l)
}

// article is draft
func (a *Article) IsDraft() bool {
	return a.Status == ARTICLE_STATUS_DRAFT
}

// article has comment
func (a *Article) HasComment() bool {
	return a.Comments > 0
}

// read article's owner
func (a *Article) User() *User {
	if a.userData == nil {
		u, err := GetUserBy("id", a.AuthorId)
		if err != nil || u == nil {
			a.userData = &User{
				Name: "Unknown",
				Nick: "Unknown",
			}
		} else {
			a.userData = u
		}
	}
	return a.userData
}

// list general articles,
// contains publish and draft articles
func ListGeneralArticle(page, size int64, order string) ([]*Article, error) {
	articles := make([]*Article, 0)
	if err := vars.Db.Where("status != ?", ARTICLE_STATUS_DELETE).
		Limit(int(size), int((page-1)*size)).OrderBy(order).Find(&articles); err != nil {
		log.Error("Db|ListGeneralArticle|%d,%d|%s|%s", page, size, order, err.Error())
		return nil, err
	}
	return articles, nil
}

// count general articles
func CountGeneralArticle() (int64, error) {
	return vars.Db.Where("status != ?", ARTICLE_STATUS_DELETE).Count(new(Article))
}

// list articles with one status
func ListStatusArticle(status string, page, size int64, order string) ([]*Article, error) {
	articles := make([]*Article, 0)
	if err := vars.Db.Where("status = ?", status).
		Limit(int(size), int((page-1)*size)).OrderBy(order).Find(&articles); err != nil {
		log.Error("Db|ListStatusArticle|%s|%d,%d|%s|%s", status, page, size, order, err.Error())
		return nil, err
	}
	return articles, nil
}

// count articles with one status
func CountStatusArticle(status string) (int64, error) {
	return vars.Db.Where("status = ?", status).Count(new(Article))
}

// get article by column and value
func GetArticleBy(col string, v interface{}) (*Article, error) {
	a := new(Article)
	if _, err := vars.Db.Where(col+" = ?", v).Get(a); err != nil {
		log.Error("Db|GetArticleBy|%s,%v|%s", col, v, err.Error())
		return nil, err
	}
	if a.Id > 0 {
		return a, nil
	}
	return nil, nil
}

// save article.
// if Article.Id, update article,
// or insert new article;
// return the saved article.
func SaveArticle(a *Article) (*Article, error) {
	if a.Id > 0 {
		if _, err := vars.Db.Where("id = ?", a.Id).
			Cols("title,link,update_time,preview,body,topic,tag_string,status,comment_status").
			Update(a); err != nil {
			log.Error("Db|SaveArticle|%d|%s", a.Id, err.Error())
			return nil, err
		}
	} else {
		if _, err := vars.Db.Insert(a); err != nil {
			log.Error("Db|SaveArticle|%d|%s", a.Id, err.Error())
			return nil, err
		}

	}
	if a.TagString != "" {
		if err := saveTags(a.Id, a.TagString); err != nil {
			return nil, err
		}
	}
	return GetArticleBy("id", a.Id)
}

func saveTags(id int64, tagStr string) error {
	// delete old tags
	if _, err := vars.Db.Where("article_id = ?", id).Delete(new(Tag)); err != nil {
		log.Error("Db|SaveTags|%d,%s|%s", id, tagStr, err.Error())
		return err
	}
	// save new tags
	tags := strings.Split(strings.Replace(tagStr, "，", ",", -1), ",")
	for _, t := range tags {
		if _, err := vars.Db.Insert(&Tag{ArticleId: id, Tag: t}); err != nil {
			log.Error("Db|SaveTags|%d,%s|%s", id, t, err.Error())
			return err
		}
	}
	return nil
}

// remove article by id
func RemoveArticle(id int64) error {
	a := new(Article)
	a.Status = ARTICLE_STATUS_DELETE
	if _, err := vars.Db.Where("id = ?", id).Cols("status").Update(a); err != nil {
		log.Error("Db|RemoveArticle|%d|%s", id, err.Error())
		return err
	}
	return nil
}

// tag struct
type Tag struct {
	Id        int64
	ArticleId int64
	Tag       string
}
