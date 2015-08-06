package base

import (
	"github.com/Unknwon/i18n"
	"github.com/lunny/tango"
	"strings"
)

var (
	langCookieName = "lang"
	langParamName  = "lang"

	_ ILanguager = new(BaseLanguager)
)

// language interface
type ILanguager interface {
	SetLang(lang string)
	LangTr(format string, args ...interface{}) string
}

// base language struct
type BaseLanguager struct {
	Lang string
}

// set language
func (bl *BaseLanguager) SetLang(lang string) {
	bl.Lang = lang
}

// current translate method
func (bl *BaseLanguager) LangTr(format string, args ...interface{}) string {
	return i18n.Tr(bl.Lang, format, args...)
}

type i18nHandle struct {
	Lang string
}

func (i *i18nHandle) Tr(format string, args ...interface{}) string {
	return i18n.Tr(i.Lang, format, args...)
}

// language middleware handler
func I18nHandler() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		// get language
		lang := getContextLanguage(ctx)

		// set to interface
		if l, ok := ctx.Action().(ILanguager); ok {
			l.SetLang(lang)
		}

		// set to view
		if r, ok := ctx.Action().(IRender); ok {
			r.Assign("i18n", &i18nHandle{lang})
		}
		ctx.Next()
	}
}

// get language from context
func getContextLanguage(ctx *tango.Context) string {
	// get from cookie
	lang := ctx.Cookie(langCookieName)

	// get from header
	if lang == "" {
		al := ctx.Req().Header.Get("Accept-Language")
		if len(al) > 4 {
			lang = al[:5] // Only compare first 5 letters.
		}
	}

	// get from query param
	if lang == "" {
		lang = ctx.Param(langParamName)
	}

	// get default if not find in context
	lang = strings.ToLower(lang)
	if !i18n.IsExist(lang) {
		lang = i18n.GetLangByIndex(0)
	}

	return lang
}
