package admin

import (
	"errors"
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/fuxiaohei/purine/utils"
	"github.com/lunny/tango"
	"strings"
)

type Media struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (m *Media) Get() {
	opt := &mapi.MediaListOption{
		Page: m.FormInt64("page", 1),
		Size: 10,
	}
	res := mapi.Call(mapi.ListMedia, opt)
	if !res.Status {
		m.RenderError(500, errors.New(res.Error))
		return
	}
	m.Assign("Media", res.Data["media"].([]*model.Media))
	m.Assign("Pager", res.Data["pager"].(*utils.Pager))
	m.Title("Media")
	m.Render("media.tmpl")
}

type Upload struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (u *Upload) Post() {
	u.Req().ParseMultipartForm(16 << 20)
	meta := &mapi.UploadMediaMeta{
		u.Ctx,
		u.AuthUser,
	}
	res := mapi.Call(mapi.UploadMedia, meta)
	u.ServeJson(res)
}

type MediaDelete struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (md *MediaDelete) Get() {
	if !strings.Contains(md.Req().Referer(), "/admin/media") {
		md.RenderError(401, nil)
		return
	}
	id := md.FormInt64("id", 0)
	if id < 1 {
		md.Redirect("/admin/media")
		return
	}
	res := mapi.Call(mapi.DelMedia, id)
	if !res.Status {
		md.RenderError(500, errors.New(res.Error))
		return
	}
	md.Redirect("/admin/media")
	return
}
