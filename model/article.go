package model

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

type Tag struct {
	Id        int64
	ArticleId int64
	Tag       string
}
