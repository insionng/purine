package model

import (
	"github.com/fuxiaohei/purine/src/log"
	"github.com/fuxiaohei/purine/src/vars"
)

const (
	COMMENT_FROM_ARTICLE = "article"
	COMMENT_FROM_PAGE    = "page"

	COMMENT_STATUS_APPROVED = "approved"
	COMMENT_STATUS_WAIT     = "wait"
	COMMENT_STATUS_SPAM     = "spam"
	COMMENT_STATUS_DELETED  = "deleted"
)

type Comment struct {
	Id         int64
	Name       string
	UserId     int64
	Email      string
	Url        string
	AvatarUrl  string
	Body       string
	CreateTime int64  `xorm:"created"`
	Status     string `xorm:"index(status)"`

	UserIp    string
	UserAgent string

	From     string `xorm:"index(from)"`
	FromId   int64  `xorm:"index(from)"`
	ParentId int64  `xorm:"index(parent)"`
}

// save a comment,
// always insert new
func SaveComment(c *Comment) error {
	if _, err := vars.Db.Insert(c); err != nil {
		log.Error("Db|SaveComment|%s", err.Error())
		return err
	}
	return nil
}

// get a comment by column and value
func GetCommentBy(col string, v interface{}) (*Comment, error) {
	c := new(Comment)
	if isIdColumn(col) {
		if _, err := vars.Db.Id(v).Get(c); err != nil {
			log.Error("Db|GetCommentBy|%s,%v|%s", col, v, err.Error())
			return nil, err
		}
	} else {
		if _, err := vars.Db.Where(col+" = ?", v).Get(c); err != nil {
			log.Error("Db|GetCommentBy|%s,%v|%s", col, v, err.Error())
			return nil, err
		}
	}
	if c.Id > 0 {
		return c, nil
	}
	return nil, nil
}

// change comment status
func ChangeCommentStatus(cid int64, status string) (*Comment, error) {
	c := &Comment{Status: status}
	if _, err := vars.Db.Id(cid).Cols("status").Update(c); err != nil {
		log.Error("Db|ChangeCommentStatus|%d,%s|%s", cid, status, err.Error())
		return nil, err
	}
	return GetCommentBy("id", cid)
}

// list all comments
func ListAllCommments(page, size int64, order string) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	if err := vars.Db.OrderBy(order).
		Where("status != ?", COMMENT_STATUS_DELETED).
		Limit(int(size), int((page-1)*size)).Find(&comments); err != nil {
		log.Error("Db|ListAllCommments|%s|%d,%d|%s|%s", "all", page, size, order, err.Error())
		return nil, err
	}
	return comments, nil
}

// list comments in specific status
func ListStatusComments(status string, page, size int64, order string) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	if err := vars.Db.Where("status = ?", status).OrderBy(order).
		Limit(int(size), int((page-1)*size)).Find(&comments); err != nil {
		log.Error("Db|ListStatusComments|%s|%d,%d|%s|%s", status, page, size, order, err.Error())
		return nil, err
	}
	return comments, nil
}

// list all comments in article
func ListAllCommentsInArticle(aid, page, size int64, order string) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	if err := vars.Db.
		OrderBy(order).
		Where("status != ? AND from = ? AND from_id = ?", COMMENT_STATUS_DELETED, COMMENT_FROM_ARTICLE, aid).
		Limit(int(size), int((page-1)*size)).Find(&comments); err != nil {
		log.Error("Db|ListAllCommentsInArticle|%s|%d,%d|%s|%s", "all", page, size, order, err.Error())
		return nil, err
	}
	return comments, nil
}

// list comments by status in article
func ListStatusCommentsInArticle(status string, aid, page, size int64, order string) ([]*Comment, error) {
	comments := make([]*Comment, 0)
	if err := vars.Db.
		Where("status = ? AND from = ? AND from_id = ?", status, COMMENT_FROM_ARTICLE, aid).
		OrderBy(order).
		Limit(int(size), int((page-1)*size)).Find(&comments); err != nil {
		log.Error("Db|ListStatusCommentsInArticle|%s|%d,%d|%s|%s", status, page, size, order, err.Error())
		return nil, err
	}
	return comments, nil
}
