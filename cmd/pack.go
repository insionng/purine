package cmd

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Unknwon/cae/zip"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/utils"
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
		log.Error("Pack | only support --src flag")
	},
}

func packSrcZip() (string, error) {
	zip.Verbose = false
	// create zip file name from time unix
	filename := time.Now().Format("20060102150405.zip")
	z, e := zip.Create(filename)
	if e != nil {
		return "", e
	}
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}
	z.AddDir("static/admin", path.Join(root, "static", "admin"))
	z.AddDir("static/default", path.Join(root, "static", "default"))
	if err = z.Flush(); err != nil {
		return "", err
	}
	if e != nil {
		return "", e
	}
	z.Close()
	return filename, nil
}

func PackSrc(ctx *cli.Context) {
	t := time.Now()
	log.Info("Pack | %-8s", "Source")

	file, err := packSrcZip()
	if err != nil {
		log.Error("Pack | %-8s | %s", "Zip", err.Error())
		return
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("Pack | %-8s | %s", "Zip", err.Error())
		return
	}
	zipWriter, err := os.OpenFile("cmd/asset.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error("Pack | %-8s | %s", "Zip", err.Error())
		return
	}
	header := `package cmd
const zipBytes="`
	zipWriter.Write([]byte(header))
	encoder := base64.NewEncoder(base64.StdEncoding, zipWriter)
	encoder.Write(bytes)
	encoder.Close()
	zipWriter.Write([]byte(`"`))
	zipWriter.Sync()
	zipWriter.Close()
	if err = os.Remove(file); err != nil {
		log.Error("Pack | %-8s | %s", "Zip", err.Error())
		return
	}
	log.Info("Pack | %-8s | %s", "Zip", utils.FriendBytesSize(int64(len(bytes))))

	log.Info("Pack | %-8s | %.1f ms", "Source", time.Since(t).Seconds()*1000)
}
