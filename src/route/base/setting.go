package base

import (
	"errors"
	"github.com/fuxiaohei/purine/src/mapi"
	"github.com/lunny/tango"
)

var _ ISetting = new(BaseSetting)

// setting handler interface
type ISetting interface {
	SetGeneral(*mapi.SettingGeneral)
	GetGeneral() *mapi.SettingGeneral
	GetGeneralByKey(key string) string
}

// base setting struct, maintains all general settings
type BaseSetting struct {
	setting *mapi.SettingGeneral
}

func (bs *BaseSetting) SetGeneral(s *mapi.SettingGeneral) {
	bs.setting = s
}

func (bs *BaseSetting) GetGeneral() *mapi.SettingGeneral {
	return bs.setting
}

func (bs *BaseSetting) GetGeneralByKey(key string) string {
	return bs.setting.Get(key)
}

// setting middleware handler
func SettingHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		setter, ok := ctx.Action().(ISetting)
		if !ok {
			ctx.Next()
			return
		}

		// read general data
		res := mapi.Call(mapi.Setting.ReadGeneral, nil)
		if !res.Status {
			panic(errors.New(res.Error))
		}
		generalSettings := res.Data["general"].(*mapi.SettingGeneral)
		setter.SetGeneral(generalSettings)

		// assign general data
		t, ok := ctx.Action().(IRender)
		if ok {
			t.Assign("General", generalSettings)
		}

		ctx.Next()
	}
}
