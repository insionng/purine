package model

import (
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/vars"
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
