package base

import (
	"github.com/fuxiaohei/purine/model"
	"github.com/tango-contrib/renders"
	"net/http"
	"path"
)

type BaseRender struct {
	renders.Renderer

	viewData    map[string]interface{}
	themePrefix string
}

func (r *BaseRender) Title(title string) {
	r.Assign("Title", title)
}

func (r *BaseRender) HasAssign(key string) bool {
	_, ok := r.viewData[key]
	return ok
}

func (r *BaseRender) Assign(key string, value interface{}) {
	if len(r.viewData) == 0 {
		r.viewData = map[string]interface{}{
			"ThemeLink": "/" + path.Join("static", r.themePrefix),
		}
	}
	r.viewData[key] = value
}

func (r *BaseRender) Render(status int, tpl string) {
	if r.themePrefix != "" {
		tpl = path.Join(r.themePrefix, tpl)
	}
	if err := r.Renderer.StatusRender(status, tpl, r.viewData); err != nil {
		panic(err)
	}
}

func (r *BaseRender) RenderError(status int, err error) {
	if !r.HasAssign("Title") {
		r.Title("Error")
	}
	r.Assign("Status", status)
	if err != nil {
		r.Assign("Error", err.Error())
	}
	r.Render(status, "error.tmpl")
}

type AdminRender struct {
	BaseRender
}

func (r *AdminRender) Title(title string) {
	if r.themePrefix == "" {
		r.themePrefix = "admin"
	}
	r.BaseRender.Title(title)
}

func (r *AdminRender) Assign(key string, value interface{}) {
	if r.themePrefix == "" {
		r.themePrefix = "admin"
	}
	r.BaseRender.Assign(key, value)
}

func (r *AdminRender) Render(tpl string) {
	if r.themePrefix == "" {
		r.themePrefix = "admin"
	}
	r.BaseRender.Render(http.StatusOK, tpl)
}

func (r *AdminRender) RenderError(status int, err error) {
	if r.themePrefix == "" {
		r.themePrefix = "admin"
	}
	r.BaseRender.RenderError(status, err)
}

type ThemeRender struct {
	BaseRender
	generalSettings map[string]string
	isFillDefault   bool
}

func (t *ThemeRender) GetSetting(key string) string {
	t.fillDefault()
	return t.generalSettings[key]
}

func (t *ThemeRender) fillDefault() {
	if t.isFillDefault {
		return
	}
	if t.themePrefix == "" {
		theme, err := model.GetCurrentTheme()
		if err != nil {
			panic(err)
		}
		t.themePrefix = theme.Directory
		t.Assign("ThemeLink", "/"+path.Join("static", t.themePrefix))
	}
	// assign general data
	generalSettings, err := model.GetSettings("title", "subtitle", "desc", "keyword")
	if err != nil {
		panic(err)
	}
	t.generalSettings = generalSettings
	if !t.HasAssign("Title") {
		t.Title(generalSettings["title"])
	}
	t.Assign("General", generalSettings)
}

func (t *ThemeRender) Render(tpl string) {
	t.fillDefault()
	t.BaseRender.Render(http.StatusOK, tpl)
}

func (t *ThemeRender) RenderError(status int, err error) {
	t.fillDefault()
	if status == 404 {
		t.Title("NOT FOUND")
	}
	if status == 401 {
		t.Title("FORBBIDEN")
	}
	if status >= 500 {
		t.Title("SERVER ERROR")
	}
	t.BaseRender.RenderError(status, err)
}
