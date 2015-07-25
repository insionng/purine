package admin

import (
	"github.com/fuxiaohei/purine/mapi"
	"github.com/fuxiaohei/purine/route/base"
	"github.com/lunny/tango"
)

type Upload struct {
	base.AdminRender
	base.BaseAuther
	tango.Ctx
}

func (u *Upload) Post() {
	u.Req().ParseMultipartForm(16 << 20)
	meta := &mapi.UploadMediaMeta{
		u.Ctx,
		u.AuthUser,
	}
	res := mapi.Call(mapi.UploadMedia, meta)
	u.ServeJson(res)
}
