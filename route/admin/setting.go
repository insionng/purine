package admin

import (
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/lunny/tango"
)

type Setting struct {
	base.AdminRender
	base.BaseAuther
	base.BaseBinder
	tango.Ctx
}

func (s *Setting) Get() {
	generalSettings, err := model.GetSettings("title", "subtitle", "desc", "keyword")
	if err != nil {
		s.RenderError(500, err)
		return
	}
	s.Assign("General", generalSettings)

	themes, err := model.GetThemes()
	if err != nil {
		s.RenderError(500, err)
		return
	}
	s.Assign("Themes", themes)

	s.Title("Setting")
	s.Render("setting.tmpl")
}

func (s *Setting) Post() {
	if s.Form("general") == "true" {
		s.postGeneral()
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
	res := mapi.Call(mapi.SaveGeneral, form)
	if !res.Status {
		s.Assign("Error", res.Error)
		s.Get()
		return
	}
	s.Assign("Success", true)
	s.Get()
}
