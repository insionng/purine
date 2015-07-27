package cmd

import (
	"encoding/base64"
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
	zipWriter, _ := os.OpenFile("cmd/asset.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	header := `package cmd
const zipBytes="`
	zipWriter.Write([]byte(header))
	encoder := base64.NewEncoder(base64.StdEncoding, zipWriter)
	encoder.Write(bytes)
	encoder.Close()
	zipWriter.Write([]byte(`"`))
	zipWriter.Sync()
	zipWriter.Close()
    os.Remove(file)
	log.Info("Pack | %-8s | %s", "Zip", utils.FriendBytesSize(int64(len(bytes))))

	log.Info("Pack | %-8s | %.1f ms", "Source", time.Since(t).Seconds()*1000)
}
