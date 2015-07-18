package base

import (
	"github.com/fuxiaohei/purine/model"
	"github.com/tango-contrib/renders"
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

func (r *BaseRender) Assign(key string, value interface{}) {
	if len(r.viewData) == 0 {
		r.viewData = map[string]interface{}{
			"ThemeLink": "/" + path.Join("static", r.themePrefix),
		}
	}
	r.viewData[key] = value
}

func (r *BaseRender) Render(tpl string) {
	if r.themePrefix != "" {
		tpl = path.Join(r.themePrefix, tpl)
	}
	if err := r.Renderer.Render(tpl, r.viewData); err != nil {
		panic(err)
	}
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
	r.BaseRender.Render(tpl)
}

func (r *AdminRender) RenderError(err error) {
	r.Title("Error")
	r.Assign("Error", err.Error())
	r.Render("error.tmpl")
}

type ThemeRender struct {
	BaseRender
}

func (t *ThemeRender) Render(tpl string) {
	if t.themePrefix == "" {
		theme, err := model.GetCurrentTheme()
		if err != nil {
			panic(err)
		}
		t.themePrefix = theme.Directory
		t.Assign("ThemeLink", "/"+path.Join("static", t.themePrefix))
	}
	t.BaseRender.Render(tpl)
}
