package vars

import "github.com/codegangsta/cli"

var Cli *cli.App

func init() {
	Cli = cli.NewApp()
	Cli.Name = NAME
	Cli.Usage = DESCRIPTION
	Cli.Version = VERSION
	Cli.Author = AUTHOR + " " + AUTHOR_EMAIL
}
