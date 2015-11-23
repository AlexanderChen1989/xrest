package utils

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestDecodeJSON(t *testing.T) {
	bodyStr := `{"name": "alex", "age": 27}`
	body := bytes.NewBuffer([]byte(bodyStr))
	r, _ := http.NewRequest("GET", "/api/test", body)

	var payload struct {
		Name string
		Age  int
	}
	resetPayload := func() {
		payload.Name = ""
		payload.Age = 0
	}

	ctx := context.Background()
	bs := NewBodyStore()

	for i := 0; i < 1000; i++ {
		resetPayload()
		var err error
		ctx, err = bs.DecodeJSON(ctx, r, &payload)
		assert.Nil(t, err)
		assert.Equal(t, payload.Name, "alex")
		assert.Equal(t, payload.Age, 27)
	}
	r.Body.Close()
}
