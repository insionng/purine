package cmd

import (
	"bytes"
	"encoding/base64"
	"github.com/Unknwon/cae/zip"
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/src/log"
	"github.com/fuxiaohei/purine/src/model"
	"github.com/fuxiaohei/purine/src/utils"
	"github.com/fuxiaohei/purine/src/vars"
	"github.com/go-xorm/xorm"
	"github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"time"
)

var installCmd cli.Command = cli.Command{
	Name:  "install",
	Usage: "install new site for first run",
	Action: func(ctx *cli.Context) {
		t := time.Now()
		// if is not new site,
		if !IsNewSite(ctx) {
			log.Info("NewSite | %-8s", "Done")
			return
		}
		NewSite(ctx)
		NewSiteData(ctx)
		NewSiteAsset(ctx)
		log.Info("NewSite | %-8s | %.1fms", "Finish", time.Since(t).Seconds()*1e3)
	},
}

// check is new site
func IsNewSite(ctx *cli.Context) bool {
	if com.IsFile(vars.CONFIG_FILE) && com.IsFile(vars.DATA_FILE) {
		return false
	}
	return true
}

// new site
func NewSite(ctx *cli.Context) {
	config := model.NewConfig()

	// encode config
	if err := model.WriteConfig(config, vars.CONFIG_FILE); err != nil {
		log.Error("NewSite | %s", err.Error())
		return
	}

	log.Info("NewSite | %-8s | %s", "Init", vars.CONFIG_FILE)
	log.Info("NewSite | %-8s | %s", "Version", config.Version)
	log.Info("NewSite | %-8s | %s:%s", "Server", config.Server.Host, config.Server.Port)
}

// new site data
func NewSiteData(ctx *cli.Context) {
	sqliteVersion, _, _ := sqlite3.Version()
	log.Info("NewSite | %-8s | %s | %s", "SQLite", sqliteVersion, vars.DATA_FILE)

	engine, err := xorm.NewEngine("sqlite3", vars.DATA_FILE)
	if err != nil {
		log.Error("NewSite | %s", err.Error())
		return
	}
	engine.SetLogger(nil) // close logger

	if err = engine.Sync2(new(model.User),
		new(model.Token),
		new(model.Article),
		new(model.Page),
		new(model.Tag),
		new(model.Setting),
		new(model.Media),
		new(model.Comment)); err != nil {
		log.Error("NewSite | %s", err.Error())
		return
	}

	log.Info("NewSite | %-8s | SyncDb | %s,%s,%s,%s,%s,%s,%s,%s", "SQLite",
		reflect.TypeOf(new(model.User)).String(),
		reflect.TypeOf(new(model.Token)).String(),
		reflect.TypeOf(new(model.Article)).String(),
		reflect.TypeOf(new(model.Page)).String(),
		reflect.TypeOf(new(model.Tag)).String(),
		reflect.TypeOf(new(model.Setting)).String(),
		reflect.TypeOf(new(model.Media)).String(),
		reflect.TypeOf(new(model.Comment)).String(),
	)

	// site init data
	NewSiteInitData(engine)

	log.Info("NewSite | %-8s | Success", "SQLite")
	engine.Close()
}

// new site init data
func NewSiteInitData(engine *xorm.Engine) {
	// default user
	user := &model.User{
		Name:      "admin",
		Nick:      "admin",
		Email:     "admin@example.com",
		Url:       "#",
		AvatarUrl: utils.GravatarLink("admin@example.com"),
		Profile:   "this is an administrator",
		Role:      model.USER_ROLE_ADMIN,
		Status:    model.USER_STATUS_ACTIVE,
	}
	user.Salt = utils.Md5String("123456789")[8:24]
	user.Password = utils.Sha256String("123456789" + user.Salt)
	if _, err := engine.Insert(user); err != nil {
		log.Error("NewSite | %s", err.Error())
		return
	}

	// default article
	article := &model.Article{
		Title:         "Welcome to Purine",
		Link:          "welcome-to-purine",
		Preview:       blogPreview,
		Body:          blogContent,
		TagString:     "blog",
		Hits:          1,
		Comments:      0,
		Status:        model.ARTICLE_STATUS_PUBLISH,
		CommentStatus: model.ARTICLE_COMMENT_OPEN,
		AuthorId:      user.Id,
	}
	if _, err := engine.Insert(article); err != nil {
		log.Error("NewSite | %s", err.Error())
		return
	}

	// default settings
	settings := make([]interface{}, 0)
	settings = append(settings, &model.Setting{"title", "Purine", 0})
	settings = append(settings, &model.Setting{"subtitle", "a simple blog engine", 0})
	settings = append(settings, &model.Setting{"desc", "a simple blog engine by golang", 0})
	settings = append(settings, &model.Setting{"keyword", "purine,blog,golang", 0})
	settings = append(settings, &model.Setting{"theme", "default", 0})
	settings = append(settings, &model.Setting{"baseurl", "http://localhost:9999/", 0})
	settings = append(settings, &model.Setting{"media_imageext", ".jpg,.jpeg,.png,.gif", 0})
	settings = append(settings, &model.Setting{"media_fileext", ".txt,.zip,.doc,.xls,.ppt,.pdf", 0})
	settings = append(settings, &model.Setting{"media_nameformat", ":hash", 0})
	settings = append(settings, &model.Setting{"media_maxsize", strconv.Itoa(2 * 1024 * 1024), 0})
	if _, err := engine.Insert(settings...); err != nil {
		log.Error("NewSite | %s", err.Error())
		return
	}
}

// new site assets
func NewSiteAsset(ctx *cli.Context) {
	zip.Verbose = false
	tmpZipFile := "tmp.zip"
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(zipBytes))
	b, _ := ioutil.ReadAll(decoder)
	if err := ioutil.WriteFile(tmpZipFile, b, os.ModePerm); err != nil {
		log.Fatal("NewSite | %s", err.Error())
	}
	z, err := zip.Open(tmpZipFile)
	if err != nil {
		log.Fatal("NewSite | %s", err.Error())
	}
	z.ExtractTo("")
	z.Close()
	decoder = nil
	os.Remove(tmpZipFile)
	log.Info("NewSite | %-8s | Success", "Asset")
}

var (
	blogContent = `Welcome to ` + "`" + `Purine` + "`" + `. Now there are some tips to tell you how to use ` + "`" + `Purine` + "`" + ` blog. If you encounter any problems, raise an issue on [Github](https://github.fom/fuxiaohei/purine) or an email to [me](mailto:fuxiaohei@vip.qq.com).

## What is Purine

` + "`" + `Purine` + "`" + ` is a dynamic blog engine by [Go](https://golang.org). You can write contents by [Markdown](http://daringfireball.net/projects/markdown/) and manage them easily. It's familiar to a simple WordPress.

### Installation

You can install ` + "`" + `Purine` + "`" + ` with distributed binary files directly. [**Download.**](#) But if in manually , read following topics.

<!--more-->

##### Requirements

Install basic things: 

` + "`" + `Go 1.3+` + "`" + `

` + "`" + `SQLite` + "`" + `

Then, download and compile source codes:

	go get github.com/fuxiaohei/purine

### Setup

First step, create a new blog with default data :

` + "`" + `purine new` + "`" + `

Now a pure blog website are created. Visit ` + "`" + `http://localhost:9999` + "`" + ` to preview it.

### Management

The admin account are created by default. You can sign in admin pages from ` + "`" + `http://localhost:9999/admin/login` + "`" + ` with default user **admin** and password **123456789**.

**Warning:  Please change your admin account name and password after setup for safety.**

### Documentation

Read more documentation about configuration, command and customization at [Github Wiki](#).`

	blogPreview = `Welcome to ` + "`" + `Purine` + "`" + `. Now there are some tips to tell you how to use ` + "`" + `Purine` + "`" + ` blog. If you encounter any problems, raise an issue on [Github](https://github.fom/fuxiaohei/purine) or an email to [me](mailto:fuxiaohei@vip.qq.com).

## What is Purine

` + "`" + `Purine` + "`" + ` is a dynamic blog engine by [Go](https://golang.org). You can write contents by [Markdown](http://daringfireball.net/projects/markdown/) and manage them easily. It's familiar to a simple WordPress.

### Installation

You can install ` + "`" + `Purine` + "`" + ` with distributed binary files directly. [**Download.**](#) But if in manually , read following topics.`
)
