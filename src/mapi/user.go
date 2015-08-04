package mapi

import (
	"errors"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/utils"
)

var (
	ERR_USER_MISSING        = errors.New("user-missing")
	ERR_USER_WRONG_PASSWORD = errors.New("user-wrong-password")

	User = new(UserApi) // user api group
)

// user api group struct
type UserApi struct{}

// user login post form
type UserLoginForm struct {
	Name     string `form:"username" binding:"Required"`
	Password string `form:"password" binding:"Required"`
	Remember int64  `form:"remember"`
}

// user login
//
//  in  : *UserLoginForm
//  out : {
//          "user":*User,
//          "token":*Token,
//        }
//
func (_ *UserApi) Login(v interface{}) *Res {
	form, ok := v.(*UserLoginForm)
	if !ok {
		return Fail(paramTypeError(new(UserLoginForm)))
	}
	// get user
	user, err := model.GetUserBy("name", form.Name)
	if err != nil {
		return Fail(err)
	}
	if user == nil {
		return Fail(ERR_USER_MISSING)
	}

	// check password
	if !user.CheckPassword(form.Password) {
		return Fail(ERR_USER_WRONG_PASSWORD)
	}

	// create token
	if form.Remember == 0 {
		form.Remember = 3600
	} else {
		form.Remember = form.Remember * 3600 * 24
	}
	token, err := model.CreateToken(user.Id, form.Remember)
	if err != nil {
		return Fail(err)
	}

	// return data
	return Success(map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

// user auth by token
//
//  in  : string
//  out : {
//          "user":*User,
//          "token":*Token,
//        }
//
func (_ *UserApi) Auth(v interface{}) *Res {
	token, ok := v.(string)
	if !ok {
		return Fail(paramTypeError(""))
	}
	// get token
	t, err := model.GetValidToken(token)
	if err != nil {
		return Fail(err)
	}
	if t == nil {
		return Fail(ERR_USER_MISSING)
	}
	// get user
	user, err := model.GetUserBy("id", t.UserId)
	if err != nil {
		return Fail(err)
	}
	if user == nil {
		return Fail(ERR_USER_MISSING)
	}

	// return data
	return Success(map[string]interface{}{
		"token": token,
		"user":  user,
	})
}

var (
	ERR_USER_EMAIL       = errors.New("user-email-exist")
	ERR_PASSWORD_CONFIRM = errors.New("password-confirm-fail")
)

// user profile post form
type UserProfileForm struct {
	User    string `form:"user" binding:"Required;AlphaDashDot"`
	Nick    string `form:"nick" binding:"Required"`
	Email   string `form:"email" binding:"Required;Email"`
	Url     string `form:"url" binding:"Url"`
	Profile string `form:"profile"`
	Id      int64
}

// update user profile
//
//  in  : *UserProfileForm
//  out : {
//          "user":*User,
//        }
//
func (_ *UserApi) UpdateProfile(v interface{}) *Res {
	form, ok := v.(*UserProfileForm)
	if !ok {
		return Fail(paramTypeError(new(UserProfileForm)))
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

// user password post form
type UserPasswordForm struct {
	User    *model.User
	Old     string `form:"old" binding:"Required"`
	New     string `form:"new" binding:"Required"`
	Confirm string `form:"confirm" binding:"Required"`
}

// update user password
//
//  in  : *UserPasswordForm
//  out : nil
//
func (_ *UserApi) UpdatePassword(v interface{}) *Res {
	form, ok := v.(*UserPasswordForm)
	if !ok {
		return Fail(paramTypeError(new(UserPasswordForm)))
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
