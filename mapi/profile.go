package mapi

import (
	"errors"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/utils"
)

var (
	ERR_USER_EMAIL       = errors.New("user-email-exist")
	ERR_PASSWORD_CONFIRM = errors.New("password-confirm-fail")
)

type ProfileForm struct {
	User    string `form:"user" binding:"Required;AlphaDashDot"`
	Nick    string `form:"nick" binding:"Required"`
	Email   string `form:"email" binding:"Required;Email"`
	Url     string `form:"url" binding:"Url"`
	Profile string `form:"profile"`
	Id      int64
}

func UpdateProfile(v interface{}) *Res {
	form, ok := v.(*ProfileForm)
	if !ok {
		return Fail(paramTypeError(new(ProfileForm)))
	}
	u := &model.User{
		Name:    form.User,
		Nick:    form.Nick,
		Email:   form.Email,
		Url:     form.Url,
		Profile: form.Profile,
		Id:      form.Id,
	}
	u.AvatarUrl = utils.GravatarLink(u.Email)
	if u.Url == "" {
		u.Url = "#"
	}

	u2, err := model.GetUserBy("email", u.Email)
	if err != nil {
		return Fail(err)
	}
	if u2 != nil && u2.Id != u.Id {
		return Fail(ERR_USER_EMAIL)
	}

	if err := model.UpdateUser(u); err != nil {
		return Fail(err)
	}
	u, err = model.GetUserBy("id", u.Id)
	if err != nil {
		return Fail(err)
	}
	return Success(map[string]interface{}{
		"user": u,
	})
}

type PasswordForm struct {
	User    *model.User
	Old     string `form:"old" binding:"Required"`
	New     string `form:"new" binding:"Required"`
	Confirm string `form:"confirm" binding:"Required"`
}

func UpdatePassword(v interface{}) *Res {
	form, ok := v.(*PasswordForm)
	if !ok {
		return Fail(paramTypeError(new(PasswordForm)))
	}
	if form.Confirm != form.New {
		return Fail(ERR_PASSWORD_CONFIRM)
	}

	if !form.User.CheckPassword(form.Old) {
		return Fail(ERR_USER_WRONG_PASSWORD)
	}

	if err := model.UpdatePassword(form.User.Id, form.New); err != nil {
		return Fail(err)
	}
	return Success(nil)
}
