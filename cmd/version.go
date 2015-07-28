package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/vars"
)

var versionCmd cli.Command = cli.Command{
	Name:  "version",
	Usage: "print version string",
	Action: func(ctx *cli.Context) {
		VersionPrint(ctx)
	},
}

// print version
func VersionPrint(ctx *cli.Context) {
	println("version", vars.VERSION)
	println("publish date", vars.VERSION_DATE)
}
