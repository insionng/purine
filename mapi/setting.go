package mapi

import (
	"encoding/json"
	"fmt"
	"github.com/fuxiaohei/purine/model"
	"strconv"
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

func struct2Map(v interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	mapData := make(map[string]interface{})
	if err = json.Unmarshal(bytes, &mapData); err != nil {
		return nil, err
	}
	return mapData, nil
}

func SaveGeneralSetting(v interface{}) *Res {
	form, ok := v.(*SettingGeneralForm)
	if !ok {
		return Fail(paramTypeError(new(SettingGeneralForm)))
	}
	mapData, err := struct2Map(form)
	if err != nil {
		return Fail(err)
	}

	for k, v := range mapData {
		if err = model.SaveSetting(strings.ToLower(k), fmt.Sprint(v), 0); err != nil {
			return Fail(err)
		}
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

type SettingMediaForm struct {
	MaxSize    int64  `form:"max_size" binding:"Required"`
	ImageExt   string `form:"image_ext"`
	FileExt    string `form:"file_ext"`
	NameFormat string `form:"form_format"`
}

type SettingMedia SettingMediaForm

func ReadMediaSetting(v interface{}) *Res {
	settings, err := model.GetSettings("media_maxsize",
		"media_imageext",
		"media_fileexit",
		"media_nameformat")
	if err != nil {
		return Fail(err)
	}
	mediaSetting := &SettingMedia{
		ImageExt:   settings["media_imageext"],
		FileExt:    settings["media_fileext"],
		NameFormat: settings["media_nameformat"],
	}
	mediaSetting.MaxSize, _ = strconv.ParseInt(settings["media_maxsize"], 10, 64)
	prepareMediaSetting(mediaSetting)
	return Success(map[string]interface{}{
		"media": mediaSetting,
	})
}

func SaveMediaSetting(v interface{}) *Res {
	form, ok := v.(*SettingMediaForm)
	if !ok {
		return Fail(paramTypeError(form))
	}

	if form.NameFormat == "" {
		form.NameFormat = ":hash"
	}

	mapData, err := struct2Map(form)
	if err != nil {
		return Fail(err)
	}

	for k, v := range mapData {
		if err = model.SaveSetting("media_"+strings.ToLower(k), fmt.Sprint(v), 0); err != nil {
			return Fail(err)
		}
	}
	return Success(nil)
}

func prepareMediaSetting(m *SettingMedia) {
	if m.MaxSize == 0 {
		m.MaxSize = 1024 * 1024 // 1M
	}
	if m.ImageExt == "" {
		m.ImageExt = ".jpg,.jpeg,.png,.gif"
	}
	if m.FileExt == "" {
		m.FileExt = ".txt,.zip,.doc,.xls,.ppt,.pdf"
	}
	if m.NameFormat == "" {
		m.NameFormat = ":hash"
	}
}
