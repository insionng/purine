package cmd

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/model"
	"github.com/fuxiaohei/purine/route/admin"
	"github.com/fuxiaohei/purine/vars"
	"github.com/lunny/tango"
	"github.com/tango-contrib/binding"
	"github.com/tango-contrib/renders"
)

var servCmd cli.Command = cli.Command{
	Name:  "serv",
	Usage: "run http server to render and show pages",
	Action: func(ctx *cli.Context) {
		// read config file
		cfg := ServeConfig(ctx)
		if cfg == nil {
			log.Error("Serv|%-8s|ReadFail", "Config")
			return
		}
		log.Error("Serv|%-8s|Read|%s", "Config", configTomlFile)

		// start server
		ServeMiddleware(ctx)
		ServeRouting(ctx)
		log.Info("Serv|%-8s|%s:%s", "Http", cfg.Server.Host, cfg.Server.Port)

		vars.Server.Run(fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port))
	},
}

func ServeConfig(ctx *cli.Context) *model.Config {
	cfg := new(model.Config)
	if _, err := toml.DecodeFile(configTomlFile, cfg); err != nil {
		log.Error("Serv|%-8s|%s", "Config", err.Error())
		return nil
	}
	return cfg
}

func ServeMiddleware(ctx *cli.Context) {
	vars.Server.Use(binding.Bind())
	vars.Server.Use(tango.Static(tango.StaticOptions{
		RootPath: "static",
		Prefix:   "/static",
	}))
	vars.Server.Use(tango.Static(tango.StaticOptions{
		RootPath: "upload",
		Prefix:   "/upload",
	}))
	vars.Server.Use(renders.New(renders.Options{
		Reload:     true,
		Directory:  "static",
		Extensions: []string{".tmpl"},
	}))
}

func ServeRouting(ctx *cli.Context) {
	adminGroup := tango.NewGroup()
	adminGroup.Any("/login", new(admin.Login))
	adminGroup.Get("/", new(admin.Index))

	vars.Server.Group("/admin", adminGroup)
}
