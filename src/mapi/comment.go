package mapi

import "github.com/fuxiaohei/purine/src/model"

var (
	Comment = new(CommentApi)
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
			comments, err = model.ListStatusCommentsInArticle(opt.Status, opt.ArticleId, opt.Page, opt.Size, opt.Order)
			if err != nil {
				return Fail(err)
			}
		}
	} else {
		if opt.ArticleId > 0 {
			comments, err = model.ListAllCommments(opt.Page, opt.Size, opt.Order)
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
