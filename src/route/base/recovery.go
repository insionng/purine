package base

import (
	"errors"
	"fmt"
	"github.com/lunny/tango"
	"net/http"
	"runtime"
)

// recovery middleware handler
func Recovery(debug bool) tango.HandlerFunc {
	return func(ctx *tango.Context) {
		defer func() {
			if e := recover(); e != nil {
				content := fmt.Sprintf("Handler crashed with error: %v", e)
				for i := 1; ; i += 1 {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					} else {
						content += "\n"
					}
					content += fmt.Sprintf("%v %v", file, line)
				}

				//ctx.Logger.Error(content)

				if !ctx.Written() {
					if !debug {
						content = http.StatusText(http.StatusInternalServerError)
					}
					ctx.Result = tango.InternalServerError(content)

					if render, ok := ctx.Action().(IRender); ok {
						render.RenderError(http.StatusInternalServerError, errors.New(content))
						return
					}

					ctx.HandleError()
				}

			}
		}()

		ctx.Next()
	}
}
