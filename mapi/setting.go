package mapi

import (
	"github.com/fuxiaohei/purine/model"
	"strings"
)

type SettingGeneralForm struct {
	Title    string `form:"title" binding:"Required"`
	Subtitle string `form:"subtitle"`
	Desc     string `form:"description"`
	Keyword  string `form:"keyword"`
}

type SettingGeneral SettingGeneralForm // alias as value usage , not form binding

func (sg *SettingGeneral) Get(key string) string {
	key = strings.ToLower(key)
	switch key {
	case "title":
		return sg.Title
	case "subtitle":
		return sg.Subtitle
	case "desc":
		return sg.Desc
	case "keyword":
		return sg.Keyword
	}
	return ""
}

func SaveGeneralSetting(v interface{}) *Res {
	form, ok := v.(*SettingGeneralForm)
	if !ok {
		return Fail(paramTypeError(new(SettingGeneralForm)))
	}

	if err := model.SaveSetting("title", form.Title, 0); err != nil {
		return Fail(err)
	}

	if err := model.SaveSetting("subtitle", form.Subtitle, 0); err != nil {
		return Fail(err)
	}

	if err := model.SaveSetting("desc", form.Desc, 0); err != nil {
		return Fail(err)
	}

	if err := model.SaveSetting("keyword", form.Keyword, 0); err != nil {
		return Fail(err)
	}
	return Success(nil)
}

func ReadGeneralSetting(v interface{}) *Res {
	generalSettings, err := model.GetSettings("title", "subtitle", "desc", "keyword")
	if err != nil {
		return Fail(err)
	}
	general := &SettingGeneral{
		Title:    generalSettings["title"],
		Subtitle: generalSettings["subtitle"],
		Desc:     generalSettings["desc"],
		Keyword:  generalSettings["keyword"],
	}
	return Success(map[string]interface{}{
		"general": general,
	})
}

/*
func ListTheme(_ interface{}) *Res {
    themes,err := model.GetThemes()
    if err != nil{
        return Fail(err)
    }
    return Success(map[string]interface{}{
        "themes":themes,
    })
}*/
