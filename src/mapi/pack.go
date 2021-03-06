package mapi

import (
	"github.com/Unknwon/cae/zip"
	"os"
	"path"
	"time"
)

var (
	Pack = new(PackApi) // pack api group
)

// pack api group struct
type PackApi struct{}

// pack option
type PackOption struct {
	File        string // output file name
	IsStaticAll bool   // include static files
	IsData      bool   // include data file
}

// pack static asset and data
//
//  in  : *PackOption
//  out : {
//          "file":string
//        }
//
func (_ *PackApi) Pack(v interface{}) *Res {
	opt, ok := v.(*PackOption)
	if !ok {
		return Fail(paramTypeError(opt))
	}
	zip.Verbose = false
	// create zip file name from time unix
	if opt.File == "" {
		opt.File = time.Now().Format("20060102150405.zip")
	}
	z, err := zip.Create(opt.File)
	if err != nil {
		return Fail(err)
	}
	root, err := os.Getwd()
	if err != nil {
		return Fail(err)
	}
	if opt.IsStaticAll {
		// pack all static files
		z.AddDir("static", path.Join(root, "static"))
	} else {
		// only pack default static files
		z.AddDir("static/admin", path.Join(root, "static", "admin"))
		z.AddDir("static/default", path.Join(root, "static", "default"))
		z.AddFile("static/sitemap.xsl", path.Join(root, "static", "sitemap.xsl"))
	}
	if opt.IsData {
		// pack data
		z.AddFile("config.toml", path.Join(root, "config.toml"))
		z.AddFile("purine.db", path.Join(root, "purine.db"))
	}
	if err = z.Flush(); err != nil {
		return Fail(err)
	}
	z.Close()
	return Success(map[string]interface{}{
		"file": opt.File,
	})
}
