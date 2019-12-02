package htcache_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/gianebao/htcache"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

func TestAPI_RunWithContext(t *testing.T) {
	var (
		assert     = assert.New(t)
		mr, err    = miniredis.Run()
		okResponse = `{"status": "OK"}`
		r          = redis.NewClient(&redis.Options{Addr: mr.Addr()})

		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, okResponse)
		}))

		a = htcache.Request{
			URL:    ts.URL,
			Expiry: 1 * time.Minute,
			Redis:  r,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}

		body     = `{"foo":"bar"}`
		response []byte
	)

	assert.NoError(err)
	defer ts.Close()

	response, err = a.Exec(body)
	assert.Equal(okResponse+"\n", string(response))
	assert.False(a.FromCache, "Data must not come from cache")

	response, err = a.Exec(body)
	assert.Equal(okResponse+"\n", string(response))
	assert.True(a.FromCache, "Data must come from cache on second attempt")
}
