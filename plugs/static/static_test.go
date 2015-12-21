package static

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/AlexanderChen1989/xrest"
	"github.com/stretchr/testify/assert"
)

func equalByteSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

var testCases = []struct {
	file   string
	exists bool
}{
	{"a.txt", true},
	{"b.txt", true},
	{"c.txt", true},
	{"d.txt", false},
}

func TestStatic(t *testing.T) {
	s := New(
		Dir("test"),
		Prefix("/public"),
	)

	for i := range testCases {
		c := testCases[i]
		pipe := xrest.NewPipeline()
		pipe.Plug(s)
		pipe.Plug(
			xrest.HandlerFunc(
				func(ctx context.Context, res http.ResponseWriter, req *http.Request) {
					if c.exists {
						t.Error("should not be here\n")
						return
					}

					t.Error("should be here\n")
				},
			),
		)

		server := httptest.NewServer(pipe.HTTPHandler())
		res, err := http.Get(server.URL + "/public/" + c.file)
		fmt.Println(err)
		assert.Nil(t, err)
		bsa, err := ioutil.ReadAll(res.Body)
		assert.Nil(t, err)
		res.Body.Close()
		if c.exists {
			bsb, err := ioutil.ReadFile("test/" + c.file)
			assert.Nil(t, err)
			assert.True(t, equalByteSlice(bsa, bsb))
		}
		server.Close()
	}
}
