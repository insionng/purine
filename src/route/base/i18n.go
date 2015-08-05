package base

import (
	"github.com/Unknwon/i18n"
	"github.com/lunny/tango"
	"strings"
)

var (
	langCookieName = "lang"
	langParamName  = "lang"
)

type ILanguager interface {
	SetLang(lang string)
	LangTr(format string, args ...interface{}) string
}

type BaseLanguager struct {
	Lang string
}

func (bl *BaseLanguager) SetLang(lang string) {
	bl.Lang = lang
}

func (bl *BaseLanguager) LangTr(format string, args ...interface{}) string {
	return i18n.Tr(bl.Lang, format, args...)
}

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
			r.Assign("Lang", lang)
			r.Assign("Tr", func(format string, args ...interface{}) string {
				return i18n.Tr(lang, format, args...)
			})
		}
		ctx.Next()
	}
}

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
