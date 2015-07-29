package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route"
	"github.com/fuxiaohei/purine/route/admin"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/fuxiaohei/purine/utils"
	"github.com/fuxiaohei/purine/vars"
	"github.com/go-xorm/xorm"
	"github.com/lunny/tango"
	"github.com/mattn/go-sqlite3"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/renders"
	"html/template"
)

var servCmd cli.Command = cli.Command{
	Name:  "server",
	Usage: "run http server to render and show pages",
	Action: func(ctx *cli.Context) {
		// read config file
		cfg := ServeConfig(ctx)
		if cfg == nil {
			log.Error("Server | %-8s | ReadFail", "Config")
			return
		}
		log.Info("Server | %-8s | Read | %s", "Config", configTomlFile)

		// start Db
		ServeDb(ctx)

		// start server
		ServeMiddleware(ctx)
		ServeRouting(ctx)
		log.Info("Server | %-8s | %s:%s", "Http", cfg.Server.Host, cfg.Server.Port)

		vars.Server.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	},
}

// load config to server
func ServeConfig(ctx *cli.Context) *model.Config {
	cfg := new(model.Config)
	if _, err := toml.DecodeFile(configTomlFile, cfg); err != nil {
		log.Error("Server | %-8s | %s", "Config", err.Error())
		return nil
	}
	return cfg
}

// connect database to server
func ServeDb(ctx *cli.Context) {
	sqliteVersion, _, _ := sqlite3.Version()
	log.Info("Server | %-8s | %s | %s", "SQLite", sqliteVersion, databaseFile)

	engine, err := xorm.NewEngine("sqlite3", databaseFile)
	if err != nil {
		log.Error("Server | %s", err.Error())
		return
	}
	engine.SetLogger(nil) // close logger
	vars.Db = engine
}

// add middleware to server
func ServeMiddleware(ctx *cli.Context) {
	vars.Server.Use(base.LoggingHandler())
	for prefix, path := range vars.StaticDirectory {
		vars.Server.Use(tango.Static(tango.StaticOptions{
			RootPath: path,
			Prefix:   prefix,
		}))
	}
	vars.Server.Use(binding.Bind())
	vars.Server.Use(base.AuthHandler())
	vars.Server.Use(base.SettingHandler())
	vars.Server.Use(renders.New(renders.Options{
		Reload:     true,
		Directory:  "static",
		Extensions: []string{".tmpl"},
		Funcs: template.FuncMap{
			"Str2HTML":             utils.Str2HTML,
			"Md2HTML":              utils.Md2Html,
			"Pager2HTML":           utils.Pager2HTML,
			"Pager2Simple":         utils.Pager2HTMLSimple,
			"TimeUnixFormat":       utils.TimeUnixFormat,
			"TimeUnixFormatFriend": utils.FriendTimeUnixFormat,
			"FriendBytesSize":      utils.FriendBytesSize,
		},
	}))
}

// add routing to server
func ServeRouting(ctx *cli.Context) {
	adminGroup := tango.NewGroup()
	adminGroup.Any("/login", new(admin.Login))
	adminGroup.Get("/logout", new(admin.Logout))
	adminGroup.Any("/profile", new(admin.Profile))
	adminGroup.Post("/password", new(admin.Password))
	adminGroup.Any("/write", new(admin.Write))
	adminGroup.Any("/delete", new(admin.Delete))
	adminGroup.Get("/article", new(admin.Article))
	adminGroup.Any("/setting", new(admin.Setting))
	adminGroup.Get("/media", new(admin.Media))
	adminGroup.Get("/media/delete", new(admin.MediaDelete))
	adminGroup.Post("/upload", new(admin.Upload))
	adminGroup.Get("/", new(admin.Index))

	vars.Server.Group("/admin", adminGroup)

	vars.Server.Get("/archive", new(route.Archive))
	vars.Server.Get("/article/page/:page", new(route.Index))
	vars.Server.Get("/article/*article.html", new(route.Article))
	vars.Server.Get("/rss.xml", new(route.Rss))
	vars.Server.Get("/sitemap.xml", new(route.Sitemap))
	vars.Server.Get("/", new(route.Index))
}
