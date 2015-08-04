package main

import (
	_ "github.com/fuxiaohei/purine/src/cmd"
	"github.com/fuxiaohei/purine/src/vars"
)

func main() {
	vars.Cli.RunAndExitOnError()
}
