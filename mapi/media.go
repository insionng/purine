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

	Media = new(MediaApi) // media api group
)

// media api group struct
type MediaApi struct{}

// media upload option
type MediaUploadOption struct {
	Ctx      tango.Ctx
	User     *model.User // media's owner
	FormName string      // form field name
	IsImage  bool        // is image type
}

// upload media
//
//  in  : *MediaUploadOption
//  out : {
//          "media":*Media
//        }
//
func (_ *MediaApi) Upload(v interface{}) *Res {
	meta, ok := v.(*MediaUploadOption)
	if !ok {
		return Fail(paramTypeError(meta))
	}
	res := Setting.ReadMedia(nil)
	if !res.Status {
		return res
	}
	setting := res.Data["media"].(*SettingMedia)
	f, h, err := meta.Ctx.Req().FormFile(meta.FormName)
	if err != nil {
		return Fail(err)
	}
    defer f.Close()

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
	if meta.IsImage {
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
	if err = meta.Ctx.SaveToFile(meta.FormName, fileName); err != nil {
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

	return Success(map[string]interface{}{
		"media": m,
	})
}

// an interface to check Size() method
type fileSizer interface {
	Size() int64
}

// get file size
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

// media list option
type MediaListOption struct {
	Page, Size int64
}

// list media data
//
//  in  : *MediaListOption
//  out : {
//          "media":[]*Media,
//          "pager":*utils.Pager
//        }
//
func (_ *MediaApi) List(v interface{}) *Res {
	opt, ok := v.(*MediaListOption)
	if !ok {
		return Fail(paramTypeError(opt))
	}
	media, err := model.ListMedia(opt.Page, opt.Size)
	if err != nil {
		return Fail(err)
	}
	count, err := model.CountMedia()
	if err != nil {
		return Fail(err)
	}
	return Success(map[string]interface{}{
		"media": media,
		"pager": utils.CreatePager(opt.Page, opt.Size, count),
	})
}

// remove media file
//
//  in  : int64
//  out : nil
//
func (_ *MediaApi) Remove(v interface{}) *Res {
	id, ok := v.(int64)
	if !ok {
		return Fail(paramTypeError(id))
	}
	media, err := model.GetMediaBy("id", id)
	if err != nil {
		return Fail(err)
	}
	// remove media file
	if err = os.Remove(media.FilePath); err != nil {
		return Fail(err)
	}
	// remove database record
	if err = model.RemoveMedia(id); err != nil {
		return Fail(err)
	}
	return Success(nil)
}
