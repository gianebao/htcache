package htcache

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

// Request represents an HTTP call to be cached
type Request struct {
	URL       string
	Method    string
	Headers   map[string]string
	Expiry    time.Duration
	Redis     *redis.Client
	FromCache bool
}

var (
	// HTTPClient defines the http client to be used when doing an HTTP request
	HTTPClient   = &http.Client{}
	
	// Verbose indicates if logs will be shown
	Verbose      = false
)

func verbose(m string) {
	if Verbose {
		log.Println(m)
	}
}

// Exec executes the API and check cache
func (a *Request) Exec(body string) ([]byte, error) {
	return a.ExecWithContext(context.Background(), body)
}

// ExecWithContext executes the API call and check cache with context
func (a *Request) ExecWithContext(ctx context.Context, body string) ([]byte, error) {
	res, err := a.execFromMemory(ctx, body)

	if a.FromCache = len(res) > 0 && err == nil; a.FromCache {
		verbose("cache hit: " + a.GetID(body))
		return res, err
	}

	verbose("cache miss: " + a.GetID(body))
	return a.execHTTPWithContext(ctx, body)
}

func (a Request) execHTTPWithContext(ctx context.Context, body string) ([]byte, error) {
	var (
		req, err = http.NewRequestWithContext(ctx, a.Method, a.URL, strings.NewReader(body))
		res      *http.Response
	)

	if err != nil {
		return []byte{}, err
	}

	for i, v := range a.Headers {
		req.Header.Add(i, v)
	}

	res, err = HTTPClient.Do(req)

	if err != nil {
		return []byte{}, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	a.Redis.Set(a.GetID(body), resBody, a.Expiry)
	return resBody, err
}

func (a Request) execFromMemory(ctx context.Context, body string) ([]byte, error) {
	val, err := a.Redis.WithContext(ctx).Get(a.GetID(body)).Result()
	return []byte(val), err
}

// GetID generates a hash ID for the request
func (a Request) GetID(body string) string {
	raw := strings.Join([]string{a.Method, a.URL, body}, ":")
	return fmt.Sprintf("%x", md5.Sum([]byte(raw)))
}
