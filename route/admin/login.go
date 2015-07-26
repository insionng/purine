package admin

import (
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/lunny/tango"
	"net/http"
	"time"
)

type Login struct {
	base.AdminRender
	base.BaseBinder
	tango.Ctx
}

func (l *Login) Get() {
	l.Title("Login")
	l.Render("login.tmpl")
}

func (l *Login) Post() {
	// bind form
	form := new(mapi.LoginForm)
	if err := l.Bind(form); err != nil {
		l.Assign("Error", err.Error())
		l.Render("login.tmpl")
		return
	}

	// call login mapi
	res := mapi.Call(mapi.User.Login, form)
	if !res.Status {
		l.Assign("Error", res.Error)
		l.Render("login.tmpl")
		return
	}

	// save token to cookie
	token := res.Data["token"].(*model.Token)
	l.Cookies().Set(&http.Cookie{
		Name:     "x-token",
		Value:    token.Token,
		Path:     "/",
		Expires:  time.Unix(token.ExpireTime, 0),
		MaxAge:   int(token.ExpireTime - time.Now().Unix()),
		HttpOnly: true,
	})
	l.Redirect("/admin")
}

type Logout struct {
	tango.Ctx
}

func (l *Logout) Get() {
	l.Cookies().Set(&http.Cookie{
		Name:     "x-token",
		Value:    "",
		Path:     "/",
		MaxAge:   -3600,
		HttpOnly: false,
	})
	l.Redirect("/admin/login")
}
