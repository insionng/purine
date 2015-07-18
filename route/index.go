package route

import "github.com/fuxiaohei/purine/route/base"

type Index struct {
	base.ThemeRender
}

func (idx *Index) Get() {
	idx.Title("Index")
	idx.Render("index.tmpl")
}
