package model

import (
	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"github.com/fuxiaohei/purine/src/log"
	"github.com/fuxiaohei/purine/src/vars"
	"io/ioutil"
	"path"
	"strings"
)

// setting struct
type Setting struct {
	Name   string
	Value  string
	UserId int64
}

// save setting with key, value and owner id
func SaveSetting(key, value string, uid int64) error {
	if _, err := vars.Db.Where("name = ?", key).Delete(new(Setting)); err != nil {
		log.Error("Db|SaveSetting|%s,%s,%d|%s", key, value, uid, err.Error())
		return err
	}
	s := &Setting{
		Name:   key,
		Value:  value,
		UserId: uid,
	}
	if _, err := vars.Db.Insert(s); err != nil {
		log.Error("Db|SaveSetting|%s,%s,%d|%s", key, value, uid, err.Error())
		return err
	}
	return nil
}

// get settings by keys
func GetSettings(keys ...string) (map[string]string, error) {
	str := `"` + strings.Join(keys, `","`) + `"`
	settings := make([]*Setting, 0)
	if err := vars.Db.Where("name IN (" + str + ")").Find(&settings); err != nil {
		log.Error("Db|GetSettings|%v|%s", keys, err.Error())
		return nil, err
	}
	m := make(map[string]string)
	for _, s := range settings {
		m[s.Name] = s.Value
	}
	return m, nil
}

// theme struct
type Theme struct {
	Name      string
	Version   string
	Directory string
	IsCurrent bool
}

// get themes in diretory
func GetThemes() ([]*Theme, error) {
	themeSetting, err := GetSettings("theme")
	if err != nil {
		return nil, err
	}
	dirs, err := ioutil.ReadDir("static")
	if err != nil {
		return nil, err
	}
	themes := make([]*Theme, 0)
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		// ignore admin and upload directory
		if d.Name() == "admin" || d.Name() == "upload" {
			continue
		}

		t := new(Theme)
		t.Directory = d.Name()
		// read theme file
		tomlFile := path.Join("static", d.Name(), "theme.toml")
		if com.IsFile(tomlFile) {
			if _, err := toml.DecodeFile(tomlFile, t); err != nil {
				log.Error("Db|GetThemes|%s|%s", tomlFile, err.Error())
				return nil, err
			}
		} else {
			continue
		}

		// fill data
		if t.Name == "" {
			t.Name = d.Name()
		}
		if t.Version == "" {
			t.Version = "0.0"
		}

		// is current
		if t.Directory == themeSetting["theme"] {
			t.IsCurrent = true
		}

		themes = append(themes, t)
	}
	return themes, nil
}

// get current theme
func GetCurrentTheme() (*Theme, error) {
	themeSetting, err := GetSettings("theme")
	if err != nil {
		return nil, err
	}

	t := new(Theme)
	t.Directory = themeSetting["theme"]
	t.IsCurrent = true
	tomlFile := path.Join("static", t.Directory, "theme.toml")
	if com.IsFile(tomlFile) {
		if _, err := toml.DecodeFile(tomlFile, t); err != nil {
			log.Error("Db|GetCurrentTheme|%s|%s", tomlFile, err.Error())
			return nil, err
		}
	}
	// fill data
	if t.Name == "" {
		t.Name = t.Directory
	}
	if t.Version == "" {
		t.Version = "0.0"
	}

	return t, nil
}
