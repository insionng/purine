package mapi

import (
	"errors"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/utils"
	"strconv"
	"strings"
)

var (
	ERR_PAGE_LINK    = errors.New("page-link-exist")
	ERR_PAGE_MISSING = errors.New("page-missing")
	ERR_PAGE_PARAM   = errors.New("page-param-fail")

	Page *PageApi = new(PageApi) // page api group
)

// PAGE api group struct
type PageApi struct{}

// page list option
type PageListOption struct {
	Status string // if status, query where only status
	Page   int64
	Size   int64
	Order  string // order string, default "id DESC"
}

// fill default page list option
func preparePageListOption(opt *PageListOption) *PageListOption {
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

// list articles
//
//  in  : *ArticleListOption
//  out : {
//          "articles":[]*Page,
//          "pager":*utils.Pager
//        }
//
func (_ *PageApi) List(v interface{}) *Res {
	opt, ok := v.(*PageListOption)
	if !ok {
		return Fail(paramTypeError(opt))
	}
	opt = preparePageListOption(opt)

	if opt.Status != "" {
		pages, err := model.ListStatusPage(opt.Status, opt.Page, opt.Size, opt.Order)
		if err != nil {
			return Fail(err)
		}
		count, err := model.CountStatusPage(opt.Status)
		if err != nil {
			return Fail(err)
		}
		return Success(map[string]interface{}{
			"pages": pages,
			"pager": utils.CreatePager(opt.Page, opt.Size, count),
		})
	}

	pages, err := model.ListGeneralPage(opt.Page, opt.Size, opt.Order)
	if err != nil {
		return Fail(err)
	}
	count, err := model.CountGeneralPage()
	if err != nil {
		return Fail(err)
	}
	return Success(map[string]interface{}{
		"pages": pages,
		"pager": utils.CreatePager(opt.Page, opt.Size, count),
	})
}

// page post form
type PageForm struct {
	Title    string `form:"title" binding:"Required"`
	Link     string `form:"link" binding:"Required;AlphaDashDot"`
	Body     string `form:"body" binding:"Required"`
	Format   string `form:"format"`
	Tag      string `form:"tag"`
	Draft    string `form:"draft"`
	Id       int64  `form:"id"`
	AuthorId int64
}

// write a page,
// if input id, update page by id,
// or save new page
//
//  in  : *PageForm
//  out : {
//          "page":*Page
//        }
//
func (_ *PageApi) Write(v interface{}) *Res {
	form, ok := v.(*PageForm)
	if !ok {
		return Fail(paramTypeError(form))
	}
	// check link
	a, err := model.GetPageBy("link", form.Link)
	if err != nil {
		return Fail(err)
	}
	if a != nil {
		if form.Id == 0 || form.Id != a.Id {
			return Fail(ERR_PAGE_LINK)
		}
	}

	// create page object
	page := &model.Page{
		Id:            form.Id,
		AuthorId:      form.AuthorId,
		Title:         form.Title,
		Link:          form.Link,
		Body:          form.Body,
		Hits:          1,
		Comments:      0,
		Status:        model.PAGE_STATUS_PUBLISH,
		CommentStatus: model.PAGE_COMMENT_OPEN,
	}

	if form.Draft != "" {
		page.Status = model.ARTICLE_STATUS_DRAFT
	}

	// save page
	page, err = model.SavePage(page)
	if err != nil {
		return Fail(err)
	}

	return Success(map[string]interface{}{
		"page": page,
	})
}

// remove a page by id
//
//  in  : int64
//  out : nil
//
func (_ *PageApi) Remove(v interface{}) *Res {
	id, ok := v.(int64)
	if !ok {
		return Fail(paramTypeError(id))
	}
	if err := model.RemovePage(id); err != nil {
		return Fail(err)
	}
	return Success(nil)
}

// page route param
type PageRouteParam struct {
	Id   int64
	Link string
}

// parse page route with path and rule
func (_ *PageApi) ParseRoute(rule string, routeRule string) (*PageRouteParam, error) {
	rules := strings.Split(rule, "/")
	paramRules := strings.Split(routeRule, "/")
	if len(rules) != len(paramRules) {
		return nil, ERR_PAGE_PARAM
	}
	p := new(PageRouteParam)
	for i, r := range rules {
		if r == "id" {
			var err error
			if p.Id, err = strconv.ParseInt(paramRules[i], 10, 64); err != nil {
				return nil, err
			}
		}
		if r == "link" {
			p.Link = paramRules[i]
		}
	}
	return p, nil
}

// get a page
//
//  in  : *PageRouteParam
//  out : {
//          "page":*Page
//        }
//
func (_ *PageApi) Get(v interface{}) *Res {
	param, ok := v.(*PageRouteParam)
	if !ok {
		return Fail(paramTypeError(param))
	}
	var (
		page *model.Page
		err  error
	)
	if param.Id > 0 {
		page, err = model.GetPageBy("id", param.Id)
		if err != nil {
			return Fail(err)
		}
	}
	if param.Link != "" {
		page, err = model.GetPageBy("link", param.Link)
		if err != nil {
			return Fail(err)
		}
	}

	// check value
	if param.Id > 0 && param.Id != page.Id {
		return Fail(ERR_PAGE_MISSING)
	}
	if param.Link != "" && param.Link != page.Link {
		return Fail(ERR_PAGE_MISSING)
	}

	return Success(map[string]interface{}{
		"page": page,
	})
}
