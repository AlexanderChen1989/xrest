package i18n

import "net/http"

func FindLangDefault(r *http.Request) (lang string) {
	// Check URL parameter
	if err := r.ParseForm(); err == nil {
		if lang = r.Form.Get("lang"); len(lang) > 0 {
			return lang
		}
	}

	// Check cookie
	if cookie, err := r.Cookie("lang"); err == nil && cookie != nil {
		return cookie.Value
	}

	// Check http header 'Accept-Language'
	if al := r.Header.Get("Accept-Language"); len(al) > 0 {
		// eg. zh-CN,zh;q=0.8,en;q=0.6,ja;q=0.4,zh-TW;q=0.2
		if len(al) > 4 {
			al = al[:5] // Only compare first 5 letters
		}
		return al
	}
	return ""
}
