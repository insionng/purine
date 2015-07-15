package utils

import "html/template"

func Str2HTML(str string) template.HTML {
	return template.HTML(str)
}
