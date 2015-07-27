package cmd

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"time"

	"errors"
	"github.com/codegangsta/cli"
	"github.com/fuxiaohei/purine/log"
	"github.com/fuxiaohei/purine/mapi"
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

		// normal pack
		t := time.Now()
		opt := &mapi.PackOption{
			IsStaticAll: true,
			IsData:      true,
		}
		res := mapi.Pack.Pack(opt)
		if !res.Status {
			log.Error("Pack | %-8s | %s", "ZipAll", res.Error)
			return
		}
		file := res.Data["file"].(string)
		if fi, err := os.Stat(file); err == nil {
			log.Info("Pack | %-8s | %s | %s ", "ZipAll", file, utils.FriendBytesSize(fi.Size()))
		} else {
			log.Info("Pack | %-8s | %s", "ZipAll", file)
		}
		log.Info("Pack | %-8s | %.1fms", "ZipAll", time.Since(t).Seconds()*1000)
	},
}

func packSrcZip() (string, error) {
	opt := &mapi.PackOption{
		IsStaticAll: false,
		IsData:      false,
	}
	res := mapi.Pack.Pack(opt)
	if !res.Status {
		return "", errors.New(res.Error)
	}
	return res.Data["file"].(string), nil
}

func PackSrc(ctx *cli.Context) {
	t := time.Now()
	log.Info("Pack | %-8s", "Source")

	file, err := packSrcZip()
	if err != nil {
		log.Error("Pack | %-8s | %s", "ZipSrc", err.Error())
		return
	}

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error("Pack | %-8s | %s", "ZipSrc", err.Error())
		return
	}
	zipWriter, err := os.OpenFile("cmd/asset.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error("Pack | %-8s | %s", "ZipSrc", err.Error())
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
		log.Error("Pack | %-8s | %s", "ZipSrc", err.Error())
		return
	}
	log.Info("Pack | %-8s | %s", "ZipSrc", utils.FriendBytesSize(int64(len(bytes))))
	log.Info("Pack | %-8s | %.1fms", "ZipSrc", time.Since(t).Seconds()*1000)
}
