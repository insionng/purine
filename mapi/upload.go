package mapi

import (
	"errors"
	"github.com/lunny/tango"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

var (
	ERR_MEDIA_LARGE      = errors.New("media-too-large")
	ERR_MEDIA_WRONG_TYPE = errors.New("media-error-type")
)

func UploadMedia(v interface{}) *Res {
	ctx, ok := v.(tango.Ctx)
	if !ok {
		return Fail(paramTypeError(ctx))
	}
	res := ReadMediaSetting(nil)
	if !res.Status {
		return res
	}
	setting := res.Data["media"].(*SettingMedia)
	f, h, err := ctx.Req().FormFile("file")
	if err != nil {
		return Fail(err)
	}
	size, err := getUploadFileSize(f)
	if err != nil {
		return Fail(err)
	}
	if size > setting.MaxSize {
		return Fail(ERR_MEDIA_LARGE)
	}
	ext := path.Ext(h.Filename)
	if !strings.Contains(setting.ImageExt, ext) {
		return Fail(ERR_MEDIA_WRONG_TYPE)
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
