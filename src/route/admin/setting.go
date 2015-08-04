package admin

import (
	"errors"
	"github.com/fuxiaohei/purine/src/mapi"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/route/base"
	"github.com/lunny/tango"
)

type Setting struct {
	base.AdminRender
	base.BaseAuther
	base.BaseBinder
	tango.Ctx
}

func (s *Setting) Get() {
	res := mapi.Call(mapi.Setting.ReadGeneral, nil)
	if !res.Status {
		s.RenderError(500, errors.New(res.Error))
		return
	}
	s.Assign("General", res.Data["general"].(*mapi.SettingGeneral))

	themes, err := model.GetThemes()
	if err != nil {
		s.RenderError(500, err)
		return
	}
	s.Assign("Themes", themes)

	res = mapi.Call(mapi.Setting.ReadMedia, nil)
	if !res.Status {
		s.RenderError(500, errors.New(res.Error))
		return
	}
	s.Assign("Media", res.Data["media"].(*mapi.SettingMedia))

	s.Title("Setting")
	s.Render("setting.tmpl")
}

func (s *Setting) Post() {
	if s.Form("general") == "true" {
		s.postGeneral()
		return
	}
	if s.Form("media") == "true" {
		s.postMedia()
		return
	}
	s.Get()
}

func (s *Setting) postGeneral() {
	form := new(mapi.SettingGeneralForm)
	if err := s.Bind(form); err != nil {
		s.Assign("Error", err.Error())
		s.Get()
		return
	}
	res := mapi.Call(mapi.Setting.SaveGeneral, form)
	if !res.Status {
		s.Assign("Error", res.Error)
		s.Get()
		return
	}
	s.Assign("Success", true)
	s.Get()
}

func (s *Setting) postMedia() {
	form := new(mapi.SettingMediaForm)
	if err := s.Bind(form); err != nil {
		s.Assign("MediaError", err.Error())
		s.Get()
		return
	}
	res := mapi.Call(mapi.Setting.SaveMedia, form)
	if !res.Status {
		s.Assign("MediaError", res.Error)
		s.Get()
		return
	}
	s.Assign("MediaSuccess", true)
	s.Redirect("/admin/setting#setting-media")
}
