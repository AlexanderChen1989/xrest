package i18n

import (
	"net/http"
	"testing"

	"github.com/AlexanderChen1989/xrest"
	"golang.org/x/net/context"
)

type testPlug struct {
	zhCN bool
	t    *testing.T
	next xrest.Handler
}

func newTestPlug(zhCN bool, t *testing.T) *testPlug {
	return &testPlug{
		zhCN: zhCN,
		t:    t,
	}
}

func (tp *testPlug) Plug(h xrest.Handler) xrest.Handler {
	tp.next = h
	return tp
}

func (tp *testPlug) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	defer tp.next.ServeHTTP(ctx, w, r)
	locale, _ := FetchLocale(ctx)
	if locale == nil {
		tp.t.Error("'locale' should not be nil.")
		return
	}
	if tp.zhCN {
		if actual := locale.Tr("base.say_hi", "gopher"); actual != "嗨 gopher" {
			tp.t.Errorf("Expect %q, actual %q\n", "嗨 gopher", actual)
		}

		if actual := locale.Tr("base.hello"); actual != "你好世界！" {
			tp.t.Errorf("Expect %q, actual %q\n", "你好世界！", actual)
		}
	} else {
		if actual := locale.Tr("base.say_hi", "gopher"); actual != "Hi gopher" {
			tp.t.Errorf("Expect %q, actual %q\n", "Hi gopher", actual)
		}

		if actual := locale.Tr("base.hello"); actual != "Hello World!" {
			tp.t.Errorf("Expect %q, actual %q\n", "Hello World!", actual)
		}
	}

}

func TestI18n(t *testing.T) {
	i18nPlug, err := New("./locales", FindLangDefault)
	if err != nil {
		t.Fatal(err)
	}

	var zhCN bool

	t.Logf("Test en-US\n")
	t.Logf("Test lang from parameter\n")
	r1, err := http.NewRequest("GET", "http://xrest.org/api/hello?lang=en-US", nil)
	if err != nil {
		t.Fatal(err)
	}
	xrest.NewPipeline().
		Plug(i18nPlug, newTestPlug(zhCN, t)).
		Handler().ServeHTTP(context.Background(), nil, r1)

	t.Logf("Test lang from cookie\n")
	r2, err := http.NewRequest("GET", "http://xrest.org/api/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	r2.AddCookie(&http.Cookie{
		Name:  "lang",
		Value: "en-US",
	})
	xrest.NewPipeline().
		Plug(i18nPlug, newTestPlug(zhCN, t)).
		Handler().ServeHTTP(context.Background(), nil, r2)

	t.Logf("Test lang from http header\n")
	r3, err := http.NewRequest("GET", "http://xrest.org/api/hello", nil)
	r3.Header.Set("Accept-Language", "en-US")
	xrest.NewPipeline().
		Plug(i18nPlug, newTestPlug(zhCN, t)).
		Handler().ServeHTTP(context.Background(), nil, r3)

	// ======================================================================

	zhCN = true
	t.Logf("Test zh-CN\n")
	t.Logf("Test lang from parameter\n")
	r4, err := http.NewRequest("GET", "http://xrest.org/api/hello?lang=zh-CN", nil)
	if err != nil {
		t.Fatal(err)
	}
	xrest.NewPipeline().
		Plug(i18nPlug, newTestPlug(zhCN, t)).
		Handler().ServeHTTP(context.Background(), nil, r4)

	t.Logf("Test lang from cookie\n")
	r5, err := http.NewRequest("GET", "http://xrest.org/api/hello", nil)
	if err != nil {
		t.Fatal(err)
	}
	r5.AddCookie(&http.Cookie{
		Name:  "lang",
		Value: "zh-CN",
	})
	xrest.NewPipeline().
		Plug(i18nPlug, newTestPlug(zhCN, t)).
		Handler().ServeHTTP(context.Background(), nil, r5)

	t.Logf("Test lang from http header\n")
	r6, err := http.NewRequest("GET", "http://xrest.org/api/hello", nil)
	r6.Header.Set("Accept-Language", "zh-CN")
	xrest.NewPipeline().
		Plug(i18nPlug, newTestPlug(zhCN, t)).
		Handler().ServeHTTP(context.Background(), nil, r6)
}
