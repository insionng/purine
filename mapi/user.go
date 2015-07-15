package mapi

import (
	"errors"
	"github.com/fuxiaohei/purine/model"
)

var (
	ERR_USER_MISSING        = errors.New("user-missing")
	ERR_USER_WRONG_PASSWORD = errors.New("user-wrong-password")
)

type LoginForm struct {
	Name     string `form:"username" binding:"Required"`
	Password string `form:"password" binding:"Required"`
	Remember int64  `form:"remember"`
}

func Login(v interface{}) *Res {
	form, ok := v.(*LoginForm)
	if !ok {
		return Fail(paramTypeError(new(LoginForm)))
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