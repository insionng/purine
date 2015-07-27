package cmd

import (
	"bytes"
	"github.com/Unknwon/cae/zip"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/utils"
	"io/ioutil"
	"os"
	"path"
	"time"
)

var packCmd cli.Command = cli.Command{
	Name:  "pack",
	Usage: "pack data to zip or source",
	Flags: []cli.Flag{
		cli.BoolFlag{"src", "pack static files to sources code", ""},
	},
	Action: func(ctx *cli.Context) {
		if ctx.Bool("src") {
			PackSrc(ctx)
			return
		}
	},
}

func PackSrc(ctx *cli.Context) {
	t := time.Now()
	log.Info("Pack | %-8s", "Source")

	// get root path
	root, err := os.Getwd()
	if err != nil {
		log.Error("Pack | %-8s | %s", "Root", err.Error())
		return
	}

	// check dest file
	destFile := path.Join(root, "cmd", "_zip.go")

	// zip admin theme
	zip.Verbose = false
	var (
		buf     bytes.Buffer
		fileBuf bytes.Buffer
	)
	log.Info("Pack | %-8s | %s", "Zip", "default")
	zipWriter := zip.New(&buf)
	zipWriter.AddDir("static/admin", path.Join(root, "static", "admin"))
	zipWriter.AddDir("static/default", path.Join(root, "static", "default"))
	if err = zipWriter.Flush(); err != nil {
		log.Error("Pack | %-8s | %s", "Zip", err.Error())
		return
	}
	bytes, err := utils.Base64EncodeBytes(buf.Bytes())
	if err != nil {
		log.Error("Pack | %-8s | %s", "Zip", err.Error())
		return
	}
	log.Info("Pack | %-8s | %s | %s", "Zip", "default", utils.FriendBytesSize(int64(len(bytes))))
	zipWriter.Close()

	// write to file
	fileBuf.Write([]byte(`package cmd
const adminBytes="`))
	fileBuf.Write(bytes)
	fileBuf.WriteString(`"`)
	fileBuf.Write([]byte("\n"))

	if err = ioutil.WriteFile(destFile, fileBuf.Bytes(), os.ModePerm); err != nil {
		log.Error("Pack | %-8s | %s", "Write", err.Error())
		return
	}

	log.Info("Pack | %-8s | %.1f ms", "Source", time.Since(t).Seconds()*1000)
}
