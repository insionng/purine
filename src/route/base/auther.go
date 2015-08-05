package base

import (
	"github.com/fuxiaohei/purine/src/mapi"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/lunny/tango"
)

var _ IAuther = new(BaseAuther) // check BaseAuther implement IAuther

// authorize interface
type IAuther interface {
	GetAuthToken(*tango.Context) string
	SetAuthUser(*model.User)
	SuccessRedirect() string
	FailRedirect() string
}

// base auth handler
type BaseAuther struct {
	AuthUser *model.User
}

// get auth token
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

// set auth user
func (a *BaseAuther) SetAuthUser(u *model.User) {
	a.AuthUser = u
}

// return success redirect url string
func (a *BaseAuther) SuccessRedirect() string {
	return ""
}

// return fail url string
func (a *BaseAuther) FailRedirect() string {
	return "/admin/login"
}

// auth middleware handler
func AuthHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		auth, ok := ctx.Action().(IAuther)
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
