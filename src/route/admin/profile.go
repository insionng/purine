package admin

import (
	"github.com/fuxiaohei/purine/src/mapi"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/route/base"
)

type Profile struct {
	base.AdminRender
	base.BaseAuther
	base.BaseBinder
}

func (p *Profile) Get() {
	p.Title("Profile")
	p.Assign("User", p.AuthUser)
	p.Render("profile.tmpl")
}

func (p *Profile) Post() {
	form := new(mapi.UserProfileForm)
	if err := p.Bind(form); err != nil {
		p.Assign("Error", err.Error())
		p.Get()
		return
	}
	form.Id = p.AuthUser.Id
	res := mapi.Call(mapi.User.UpdateProfile, form)
	if !res.Status {
		p.Assign("Error", res.Error)
	} else {
		p.SetAuthUser(res.Data["user"].(*model.User))
		p.Assign("Success", true)
	}
	p.Get()
}

type Password struct {
	base.AdminRender
	base.BaseAuther
	base.BaseBinder
}

func (p *Password) Post() {
	form := new(mapi.UserPasswordForm)
	if err := p.Bind(form); err != nil {
		p.Assign("PasswordError", err.Error())
		p.Title("Profile")
		p.Assign("User", p.AuthUser)
		p.Render("profile.tmpl")
		return
	}
	form.User = p.AuthUser

	res := mapi.Call(mapi.User.UpdatePassword, form)
	if !res.Status {
		p.Assign("PasswordError", res.Error)
	} else {
		p.Assign("PasswordSuccess", true)
	}
	p.Title("Profile")
	p.Assign("User", p.AuthUser)
	p.Render("profile.tmpl")
}
