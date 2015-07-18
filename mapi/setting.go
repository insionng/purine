package mapi

import "github.com/fuxiaohei/purine/model"

type SettingGeneralForm struct {
	Title    string `form:"title" binding:"Required"`
	Subtitle string `form:"subtitle"`
	Desc     string `form:"description"`
	Keyword  string `form:"keyword"`
}

func SaveGeneral(v interface{}) *Res {
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
