package mapi

import "github.com/lunny/tango"

func UploadMedia(v interface{}) *Res {
	ctx, ok := v.(tango.Ctx)
	if !ok {
		return Fail(paramTypeError(ctx))
	}
	return Success(nil)
}
