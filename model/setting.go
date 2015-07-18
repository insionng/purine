package model

import (
	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/vars"
	"io/ioutil"
	"path"
	"strings"
)

type Setting struct {
	Name   string
	Value  string
	UserId int64
}

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

type Theme struct {
	Name      string
	Version   string
	Directory string
	IsCurrent bool
}

func GetThemes() ([]*Theme, error) {
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
		themes = append(themes, t)
	}
	return themes, nil
}
