package main

import (
	_ "github.com/fuxiaohei/purine/cmd"
	"github.com/fuxiaohei/purine/vars"
)

func main() {
	vars.Cli.RunAndExitOnError()
}
