// Package cmd contains all commands
package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/vars"
)

func init() {
	vars.Cli.Commands = []cli.Command{
		installCmd,
		versionCmd,
		servCmd,
		packCmd,
	}
}
