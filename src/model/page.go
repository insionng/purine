package model

import (
	"fmt"
	"github.com/fuxiaohei/purine/src/log"
	"github.com/fuxiaohei/purine/src/vars"
	"net/url"
	"time"
)

const (
	PAGE_STATUS_PUBLISH = "publish"
	PAGE_STATUS_DRAFT   = "draft"
	PAGE_STATUS_DELETE  = "deleted"

	PAGE_COMMENT_OPEN  = "open"
	PAGE_COMMENT_CLOSE = "closed"
	PAGE_COMMENT_WAIT  = "wait" // waiting 30 days to close
)

// page struct
type Page struct {
	Id       int64
	AuthorId int64

	Title      string
	CreateTime int64  `xorm:"created"`
	UpdateTime int64  `xorm:"updated"`
	Link       string `xorm:"unique"`
	Body       string
	Topic      string
	Template   string

	Hits          int64
	Comments      int64
	Status        string
	CommentStatus string

	userData *User `xorm:"-"`
}

// format page create time
func (p *Page) Date(layout string) string {
	return time.Unix(p.CreateTime, 0).Format(layout)
}

// page visitable link
func (p *Page) Href() string {
	l := p.Link
	if l == "" {
		l = url.QueryEscape(p.Title)
	}
	return fmt.Sprintf("%s.html", l)
}

// page is draft
func (p *Page) IsDraft() bool {
	return p.Status == ARTICLE_STATUS_DRAFT
}

// page has comment
func (p *Page) HasComment() bool {
	return p.Comments > 0
}

// read page's owner
func (p *Page) User() *User {
	if p.userData == nil {
		u, err := GetUserBy("id", p.AuthorId)
		if err != nil || u == nil {
			p.userData = &User{
				Name: "Unknown",
				Nick: "Unknown",
			}
		} else {
			p.userData = u
		}
	}
	return p.userData
}

// get page by column and value
func GetPageBy(col string, v interface{}) (*Page, error) {
	p := new(Page)
	if isIdColumn(col) {
		if _, err := vars.Db.Id(v).Get(p); err != nil {
			log.Error("Db|GetPageBy|%s,%v|%s", col, v, err.Error())
			return nil, err
		}
	} else {
		if _, err := vars.Db.Where(col+" = ?", v).Get(p); err != nil {
			log.Error("Db|GetPageBy|%s,%v|%s", col, v, err.Error())
			return nil, err
		}
	}
	if p.Id > 0 {
		return p, nil
	}
	return nil, nil
}

// save page.
// if page.Id, update page,
// or insert new page;
// return the saved page.
func SavePage(p *Page) (*Page, error) {
	if p.Id > 0 {
		if _, err := vars.Db.Where("id = ?", p.Id).
			Cols("title,link,update_time,body,topic,status,comment_status").
			Update(p); err != nil {
			log.Error("Db|SavePage|%d|%s", p.Id, err.Error())
			return nil, err
		}
	} else {
		if _, err := vars.Db.Insert(p); err != nil {
			log.Error("Db|SavePage|%d|%s", p.Id, err.Error())
			return nil, err
		}

	}
	return GetPageBy("id", p.Id)
}

// list general pages,
// contains publish and draft pages
func ListGeneralPage(page, size int64, order string) ([]*Page, error) {
	pages := make([]*Page, 0)
	if err := vars.Db.Where("status != ?", PAGE_STATUS_DELETE).
		Limit(int(size), int((page-1)*size)).OrderBy(order).Find(&pages); err != nil {
		log.Error("Db|ListGeneralPage|%d,%d|%s|%s", page, size, order, err.Error())
		return nil, err
	}
	return pages, nil
}

// count general pages
func CountGeneralPage() (int64, error) {
	return vars.Db.Where("status != ?", PAGE_STATUS_DELETE).Count(new(Page))
}

// list pages with one status
func ListStatusPage(status string, page, size int64, order string) ([]*Page, error) {
	pages := make([]*Page, 0)
	if err := vars.Db.Where("status = ?", status).
		Limit(int(size), int((page-1)*size)).OrderBy(order).Find(&pages); err != nil {
		log.Error("Db|ListStatusPage|%s|%d,%d|%s|%s", status, page, size, order, err.Error())
		return nil, err
	}
	return pages, nil
}

// count pages with one status
func CountStatusPage(status string) (int64, error) {
	return vars.Db.Where("status = ?", status).Count(new(Page))
}

// remove page by id
func RemovePage(id int64) error {
	a := new(Page)
	a.Status = PAGE_STATUS_DELETE
	if _, err := vars.Db.Where("id = ?", id).Cols("status").Update(a); err != nil {
		log.Error("Db|RemovePage|%d|%s", id, err.Error())
		return err
	}
	return nil
}
