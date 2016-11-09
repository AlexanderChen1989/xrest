package i18n

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlexanderChen1989/xrest"
	"github.com/beego/i18n"
	"golang.org/x/net/context"
)

var ctxI18nKey uint8

func FetchLocale(ctx context.Context) (*i18n.Locale, bool) {
	l, ok := ctx.Value(&ctxI18nKey).(*i18n.Locale)
	return l, ok
}

type FindLangFunc func(r *http.Request) (lang string)

type I18nPlug struct {
	next     xrest.Handler
	findLang FindLangFunc
	locales  []*i18n.Locale
}

func New(localeDir string, f FindLangFunc) (p *I18nPlug, err error) {
	/*
		/etc/xrest/locale/zh-CN.ini
		/etc/xrest/locale/en-US.ini
	*/
	var locales []*i18n.Locale
	if err = filepath.Walk(localeDir, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		ext := filepath.Ext(path)
		if info.IsDir() || ext != ".ini" {
			return nil
		}

		lang := strings.ToLower(strings.TrimRight(info.Name(), ext))
		if smErr := i18n.SetMessage(lang, path); smErr != nil {
			return smErr
		}
		locales = append(locales, &i18n.Locale{Lang: lang})
		return nil
	}); err != nil {
		return nil, err
	}

	return &I18nPlug{
		findLang: f,
		locales:  locales,
	}, nil
}

func (p *I18nPlug) Plug(h xrest.Handler) xrest.Handler {
	p.next = h
	return p
}

func (p *I18nPlug) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	lang := p.findLang(r)
	ctx = context.WithValue(ctx, &ctxI18nKey, p.findLocale(lang))
	p.next.ServeHTTP(ctx, w, r)
}

func (p *I18nPlug) findLocale(lang string) *i18n.Locale {
	lang = strings.ToLower(lang)
	for i := range p.locales {
		if p.locales[i].Lang == lang {
			return p.locales[i]
		}
	}
	return p.locales[0]
}
