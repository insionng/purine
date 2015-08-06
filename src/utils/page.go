package utils

import (
	"fmt"
	"github.com/Unknwon/i18n"
	"html/template"
	"strconv"
)

// pager struct
type Pager struct {
	Current int64
	All     int64
	Pages   int64
	Size    int64
}

// create pager
func CreatePager(page, size, all int64) *Pager {
	p := &Pager{
		Current: page,
		Size:    size,
		All:     all,
	}
	p.Pages = all / size
	if all%size > 0 {
		p.Pages++
	}
	return p
}

// pager to HTML with number elements
func Pager2HTML(p *Pager, layout, lang string) template.HTML {
	tpl := ` <ul class="pager">`
	for i := 1; i <= int(p.Pages); i++ {
		if i == int(p.Current) {
			tpl += `<li><a class="current" href="` + fmt.Sprintf(layout, i) + `">` + strconv.Itoa(i) + `</a></li>`
		} else {
			tpl += `<li><a href="` + fmt.Sprintf(layout, i) + `">` + strconv.Itoa(i) + `</a></li>`
		}
	}
	tpl += "</ul>"
	return template.HTML(tpl)
}

// pager to HTML with navigator elements
func Pager2HTMLSimple(p *Pager, layout, lang string) template.HTML {
	tpl := `<div class="pager clear">`
	if p.Current > 1 {
		tpl += `<a class="prev left" href="` + fmt.Sprintf(layout, p.Current-1) + `">` + i18n.Tr(lang, "pager.prev") + `</a>`
	}
	if p.Current < p.Pages {
		tpl += `<a class="next right" href="` + fmt.Sprintf(layout, p.Current+1) + `">` + i18n.Tr(lang, "pager.next") + `</a>`
	}
	tpl += "</div>"
	return template.HTML(tpl)
}
