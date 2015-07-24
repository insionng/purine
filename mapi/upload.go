package mapi

import (
	"fmt"
	"github.com/lunny/tango"
)

func UploadMedia(v interface{}) *Res {
	ctx, ok := v.(tango.Ctx)
	if !ok {
		return Fail(paramTypeError(ctx))
	}
	_, h, e := ctx.Req().FormFile("file")
	fmt.Println(h.Filename, h.Header, e)
	return Success(nil)
}
