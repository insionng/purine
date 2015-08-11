package mapi

import (
	"errors"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/utils"
	"strings"
)

var (
	Comment = new(CommentApi)

	ERR_COMMENT_MISSING_FROM   = errors.New("missing-from")
	ERR_COMMENT_TOO_MANY_LINKS = errors.New("too-many-links")
	ERR_COMMENT_TOO_LONG       = errors.New("too-long")
	ERR_COMMENT_EMAIL_STRING   = errors.New("email-string")
)

type CommentApi struct{}

type CommentListOption struct {
	Page, Size int64
	Order      string
	Status     string
	ArticleId  int64
}

func prepareCommentListOption(opt *CommentListOption) *CommentListOption {
	if opt.Size == 0 {
		opt.Size = 10
	}
	if opt.Page == 0 {
		opt.Page = 1
	}
	if opt.Order == "" {
		opt.Order = "id ASC"
	}
	return opt
}

func (_ *CommentApi) List(v interface{}) *Res {
	opt, ok := v.(*CommentListOption)
	if !ok {
		return Fail(paramTypeError(opt))
	}
	opt = prepareCommentListOption(opt)
	var (
		comments []*model.Comment
		err      error
	)
	if opt.Status == "" {
		if opt.ArticleId > 0 {
			comments, err = model.ListAllCommentsInArticle(opt.ArticleId, opt.Page, opt.Size, opt.Order)
			if err != nil {
				return Fail(err)
			}
		} else {
			comments, err = model.ListAllComments(opt.Page, opt.Size, opt.Order)
			if err != nil {
				return Fail(err)
			}
		}
	} else {
		if opt.ArticleId > 0 {
			comments, err = model.ListStatusCommentsInArticle(opt.Status, opt.ArticleId, opt.Page, opt.Size, opt.Order)
			if err != nil {
				return Fail(err)
			}
		} else {
			comments, err = model.ListStatusComments(opt.Status, opt.Page, opt.Size, opt.Order)
			if err != nil {
				return Fail(err)
			}
		}
	}

	return Success(map[string]interface{}{
		"comments": comments,
		"article":  opt.ArticleId,
	})
}

type CommentForm struct {
	Name   string `form:"name" binding:"Required"`
	Email  string `form:"email" binding:"Required;Email"`
	Url    string `form:"url" binding:"Url"`
	Body   string `form:"body" binding:"Required"`
	Parent int64  `form:"parent"`
	For    string `form:"-"`
	ForId  int64  `form:"-"`
	UserId int64  `form:"-"`
}

var (
	emails = []string{
		"@qq.com",
		"@163.com",
		"@sina.com",
		"@126.com",
		"@gmail.com",
		"@outlook.com",
	}
)

func (_ *CommentApi) Filter(v interface{}) *Res {
	form, ok := v.(*CommentForm)
	if !ok {
		return Fail(paramTypeError(form))
	}

	// check source
	if form.For == "" || form.ForId == 0 {
		return Fail(ERR_COMMENT_MISSING_FROM)
	}

	// too long
	if len(form.Body) > 255 {
		return Fail(ERR_COMMENT_TOO_LONG)
	}
	body := strings.Replace(form.Body, " ", "", -1) // clean black spaces

	// check link
	if strings.Count(body, "href") > 0 {
		return Fail(ERR_COMMENT_TOO_MANY_LINKS)
	}

	// check email
	for _, e := range emails {
		if strings.Contains(body, e) {
			return Fail(ERR_COMMENT_EMAIL_STRING)
		}
	}

	return Success(map[string]interface{}{
		"form": form,
	})
}

func (_ *CommentApi) Save(v interface{}) *Res {
	form, ok := v.(*CommentForm)
	if !ok {
		return Fail(paramTypeError(form))
	}

	c := &model.Comment{
		Name:      form.Name,
		UserId:    form.UserId,
		Email:     form.Email,
		Url:       form.Url,
		AvatarUrl: utils.GravatarLink(form.Email),
		Body:      utils.Nl2BrString(form.Body),
		From:      form.For,
		FromId:    form.ForId,
		ParentId:  form.Parent,
		Status:    model.COMMENT_STATUS_WAIT,
	}
	if model.CountApprovedCommentsByEmail(form.Email) > 0 {
		c.Status = model.COMMENT_STATUS_APPROVED
	}
	if err := model.SaveComment(c); err != nil {
		return Fail(err)
	}
	return Success(map[string]interface{}{
		"comment": c,
	})
}
