package limit

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlexanderChen1989/xrest"
	"github.com/stretchr/testify/assert"
)

func TestLimiter(t *testing.T) {
	limiter := New(1, time.Second)

	pipe := xrest.NewPipeline()

	pipe.Plug(limiter)

	ts := httptest.NewServer(pipe.HTTPHandler())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	for i := 0; i < 10; i++ {
		res, err = http.Get(ts.URL)
		assert.Nil(t, err)
		assert.NotEmpty(t, res.StatusCode, http.StatusOK)
		// body, err := ioutil.ReadAll(res.Body)
		//
		// fmt.Println(">>>>", err, string(body))
	}
}
