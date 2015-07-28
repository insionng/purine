package model

import (
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/vars"
)

const (
	MEDIA_TYPE_IMAGE = "image"
	MEDIA_TYPE_FILE  = "file"
)

// media struct
type Media struct {
	Id int64

	Name     string
	FileName string
	FilePath string `xorm:"not null"`
	FileSize int64
	FileType string `xorm:"not null"`
	OwnerId  int64  `xorm:"not null"`

	CreateTime int64 `xorm:"created"`
	Downloads  int
}

// save new media
func SaveMedia(m *Media) error {
	if _, err := vars.Db.Insert(m); err != nil {
		return err
	}
	return nil
}

// list media files
func ListMedia(page, size int64) ([]*Media, error) {
	media := make([]*Media, 0)
	if err := vars.Db.
		Limit(int(size), int((page-1)*size)).
		OrderBy("id DESC").
		Find(&media); err != nil {
		log.Error("Db|ListMedia|%d,%d|%s", page, size, err.Error())
		return nil, err
	}
	return media, nil
}

// count media files
func CountMedia() (int64, error) {
	return vars.Db.Count(new(Media))
}

// get media by column and value
func GetMediaBy(col string, value interface{}) (*Media, error) {
	m := new(Media)
	if _, err := vars.Db.Where(col+" = ?", value).Get(m); err != nil {
		return nil, err
	}
	return m, nil
}

// remove media by id
func RemoveMedia(id int64) error {
	_, err := vars.Db.Where("id = ?", id).Delete(new(Media))
	return err
}
