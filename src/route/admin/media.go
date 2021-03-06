package admin

import (
	"errors"
	"github.com/fuxiaohei/purine/src/mapi"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/route/base"
	"github.com/fuxiaohei/purine/src/utils"
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
	res := mapi.Call(mapi.Media.List, opt)
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
	meta := &mapi.MediaUploadOption{
		Ctx:      u.Ctx,
		User:     u.AuthUser,
		FormName: "file",
		IsImage:  false,
	}
	if u.Form("type") == "image" {
		meta.IsImage = true
	}
	if u.Form("from") == "editor" {
		meta.FormName = "editormd-image-file"
	}
	res := mapi.Call(mapi.Media.Upload, meta)
	if u.Form("from") == "editor" {
		m := make(map[string]interface{})
		m["success"] = 1
		if !res.Status {
			m["success"] = 0
		}
		m["message"] = res.Error
		m["url"] = ""
		if media, ok := res.Data["media"].(*model.Media); ok {
			m["url"] = "/" + media.FilePath
		}
		u.ServeJson(m)
		return
	}
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
	res := mapi.Call(mapi.Media.Remove, id)
	if !res.Status {
		md.RenderError(500, errors.New(res.Error))
		return
	}
	md.Redirect("/admin/media")
	return
}
