package base

import (
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/model"
	"github.com/lunny/tango"
)

var _ Auther = new(BaseAuther)

type Auther interface {
	GetAuthToken(*tango.Context) string
	SetAuthUser(*model.User)
	SuccessRedirect() string
	FailRedirect() string
}

type BaseAuther struct {
	AuthUser *model.User
}

func (a *BaseAuther) GetAuthToken(ctx *tango.Context) string {
	var token string
	if token = ctx.Header().Get("X-Token"); token != "" {
		return token
	}
	if token = ctx.Cookie("x-token"); token != "" {
		return token
	}
	return ctx.Form("x-token")
}

func (a *BaseAuther) SetAuthUser(u *model.User) {
	a.AuthUser = u
}

func (a *BaseAuther) SuccessRedirect() string {
	return ""
}

func (a *BaseAuther) FailRedirect() string {
	return "/admin/login"
}

func AuthHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		auth, ok := ctx.Action().(Auther)
		if !ok {
			ctx.Next()
			return
		}
		// read token
		token := auth.GetAuthToken(ctx)
		if token != "" {
			result := mapi.Call(mapi.User.Auth, token)
			if result.Status {
				auth.SetAuthUser(result.Data["user"].(*model.User))
				ctx.Next()
				return
			}
		}
		// fail redirect
		if ctx.Req().Method == "GET" {
			if url := auth.FailRedirect(); url != "" {
				ctx.Redirect(url, 302)
				return
			}
		}

		// auth fail , no redirect, to show 401
		ctx.WriteHeader(401)
	}
}
