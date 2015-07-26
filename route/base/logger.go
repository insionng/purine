package base

import (
	"fmt"
	"github.com/fuxiaohei/purine/vars"
	"github.com/lunny/tango"
	"strings"
	"time"
)

var (
	logFormat      = "Http   | %s | %s | %s | %s "
	logErrorFormat = "Http   | %s | %s | %s | %s | %v"
)

func LoggingHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		start := time.Now()
		p := ctx.Req().URL.Path
		if len(ctx.Req().URL.RawQuery) > 0 {
			p = p + "?" + ctx.Req().URL.RawQuery
		}

		if action := ctx.Action(); action != nil {
			if l, ok := action.(tango.LogInterface); ok {
				l.SetLogger(ctx.Logger)
			}
		}

		ctx.Next()

		if !ctx.Written() {
			if ctx.Result == nil {
				ctx.Result = tango.NotFound()
			}
			ctx.HandleError()
		}

		// skip static files
		if ctx.Req().Method == "GET" {
			for prefix, _ := range vars.StaticDirectory {
				if strings.HasPrefix(p, prefix) {
					return
				}
			}
		}

		statusCode := ctx.Status()

		if statusCode >= 200 && statusCode < 400 {
			ctx.Info(
				fmt.Sprintf(
					logFormat,
					friendRemoteString(ctx.Req().RemoteAddr),
					ctx.Req().Method,
					p,
					time.Since(start)))
		} else {
			ctx.Error(
				fmt.Sprintf(logErrorFormat,
					friendRemoteString(ctx.Req().RemoteAddr),
					ctx.Req().Method,
					p,
					time.Since(start),
					ctx.Result))
		}
	}
}

func friendRemoteString(remote string) string {
	if len(remote) < 21 {
		remote += strings.Repeat(" ", 21-len(remote))
	}
	return remote
}