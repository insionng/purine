package model

import "github.com/fuxiaohei/purine/vars"

const (
	MEDIA_TYPE_IMAGE = "image"
	MEDIA_TYPE_FILE  = "file"
)

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

func SaveMedia(m *Media) error {
	if _, err := vars.Db.Insert(m); err != nil {
		return err
	}
	return nil
}
