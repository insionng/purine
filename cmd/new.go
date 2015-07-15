package cmd

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"github.com/Unknwon/com"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/model"
	"github.com/go-xorm/xorm"
	"io/ioutil"
	"os"

	"github.com/fuxiaohei/purine/utils"
	"github.com/mattn/go-sqlite3"
	"reflect"
	"time"
)

var newCmd cli.Command = cli.Command{
	Name:  "new",
	Usage: "create new site for first run",
	Action: func(ctx *cli.Context) {
		t := time.Now()
		// if is not new site,
		if !IsNewSite(ctx) {
			log.Info("NewSite|%-8s", "Done")
			return
		}
		NewSite(ctx)
		NewSiteData(ctx)
		log.Info("NewSite|%-8s|%.1fms", "Finish", time.Since(t).Seconds()*1e3)
	},
}

var (
	configTomlFile = "config.toml"
	databaseFile   = "purine.db"
)

// new site
func NewSite(ctx *cli.Context) {
	config := model.NewConfig()

	// encode config
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	if err := encoder.Encode(config); err != nil {
		log.Error("NewSite|%s", err.Error())
		return
	}

	// write config.toml
	if err := ioutil.WriteFile(configTomlFile, buf.Bytes(), os.ModePerm); err != nil {
		log.Error("NewSite|%s", err.Error())
		return
	}

	log.Info("NewSite|%-8s|%s", "Init", configTomlFile)
	log.Info("NewSite|%-8s|%s", "Version", config.Version)
	log.Info("NewSite|%-8s|%s:%s", "Server", config.Server.Host, config.Server.Port)
}

// new site data
func NewSiteData(ctx *cli.Context) {
	sqliteVersion, _, _ := sqlite3.Version()
	log.Info("NewSite|%-8s|%s|%s", "SQLite", sqliteVersion, databaseFile)

	engine, err := xorm.NewEngine("sqlite3", databaseFile)
	if err != nil {
		log.Error("NewSite|%s", err.Error())
		return
	}
	engine.SetLogger(nil) // close logger

	if err = engine.Sync2(new(model.User),
		new(model.Token),
		new(model.Article),
		new(model.Tag)); err != nil {
		log.Error("NewSite|%s", err.Error())
		return
	}

	log.Info("NewSite|%-8s|SyncDb|%s,%s,%s,%s", "SQLite",
		reflect.TypeOf(new(model.User)).String(),
		reflect.TypeOf(new(model.Token)).String(),
		reflect.TypeOf(new(model.Article)).String(),
		reflect.TypeOf(new(model.Tag)).String(),
	)

	// site init data
	NewSiteInitData(engine)

	log.Info("NewSite|%-8s|Success", "SQLite")
	engine.Close()
}

// new site init data
func NewSiteInitData(engine *xorm.Engine) {
	// default user
	user := &model.User{
		Name:      "admin",
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
		log.Error("NewSite|%s", err.Error())
		return
	}

	// default article
	article := &model.Article{
		Title:   "简洁的静态博客构建工具 —— 纸小墨（InkPaper）",
		Link:    "blog-first",
		Preview: "纸小墨（InkPaper）是一个使用GO语言编写的静态博客构建工具，可以快速搭建博客网站。优点是无依赖跨平台，配置简单构建快速，注重于简洁易用与排版优化。",
		Body: `# 纸小墨简介
        纸小墨（InkPaper）是一个使用GO语言编写的静态博客构建工具，可以快速搭建博客网站。优点是无依赖跨平台，配置简单构建快速，注重于简洁易用与排版优化。`,
		TagString:     "blog",
		Hits:          1,
		Comments:      0,
		Status:        model.ARTICLE_STATUS_PUBLISH,
		CommentStatus: model.ARTICLE_COMMENT_OPEN,
		AuthorId:      user.Id,
	}
	if _, err := engine.Insert(article); err != nil {
		log.Error("NewSite|%s", err.Error())
		return
	}
}

// check is new
func IsNewSite(ctx *cli.Context) bool {
	if com.IsFile(configTomlFile) && com.IsFile(databaseFile) {
		return false
	}
	return true
}
