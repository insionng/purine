package mapi

import (
	"errors"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/utils"
	"github.com/lunny/tango"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
)

var (
	ERR_MEDIA_LARGE      = errors.New("media-too-large")
	ERR_MEDIA_WRONG_TYPE = errors.New("media-error-type")
)

type UploadMediaMeta struct {
	Ctx  tango.Ctx
	User *model.User
}

func UploadMedia(v interface{}) *Res {
	meta, ok := v.(*UploadMediaMeta)
	if !ok {
		return Fail(paramTypeError(meta))
	}
	res := ReadMediaSetting(nil)
	if !res.Status {
		return res
	}
	setting := res.Data["media"].(*SettingMedia)
	f, h, err := meta.Ctx.Req().FormFile("file")
	if err != nil {
		return Fail(err)
	}

	// check file size
	size, err := getUploadFileSize(f)
	if err != nil {
		return Fail(err)
	}
	if size > setting.MaxSize {
		return Fail(ERR_MEDIA_LARGE)
	}

	// check ext
	ext := path.Ext(h.Filename)
	extRule := setting.FileExt
	if meta.Ctx.Form("type") == "image" {
		extRule = setting.ImageExt
	}
	if !strings.Contains(extRule, ext) {
		return Fail(ERR_MEDIA_WRONG_TYPE)
	}

	// hash file name, make dir
	now := time.Now()
	hashName := utils.Md5String(fmt.Sprintf("%d%s%d", meta.User.Id, h.Filename, now.UnixNano())) + ext
	fileName := path.Join("static/upload", hashName)
	fileDir := path.Dir(fileName)
	if !com.IsDir(fileDir) {
		if err = os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return Fail(err)
		}
	}
	if err = meta.Ctx.SaveToFile("file", fileName); err != nil {
		return Fail(err)
	}

	// save file media info
	m := &model.Media{
		Name:     h.Filename,
		FileName: hashName,
		FilePath: fileName,
		FileSize: size,
		FileType: h.Header.Get("Content-Type"),
		OwnerId:  meta.User.Id,
	}
	if err = model.SaveMedia(m); err != nil {
		return Fail(err)
	}

	return Success(nil)
}

type fileSizer interface {
	Size() int64
}

func getUploadFileSize(f multipart.File) (int64, error) {
	// if return *http.sectionReader, it is alias to *io.SectionReader
	if s, ok := f.(fileSizer); ok {
		return s.Size(), nil
	}
	// or *os.File
	if fp, ok := f.(*os.File); ok {
		fi, err := fp.Stat()
		if err != nil {
			return 0, err
		}
		return fi.Size(), nil
	}
	return 0, nil
}
