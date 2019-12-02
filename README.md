# htcache
--
    import "github.com/gianebao/htcache"


## Usage

#### type Request

```go
type Request struct {
	URL       string
	Method    string
	Headers   map[string]string
	Expiry    time.Duration
	Redis     *redis.Client
	FromCache bool
}
```

Request represents an HTTP call to be cached

#### func (*Request) Exec

```go
func (a *Request) Exec(body string) ([]byte, error)
```
Exec executes the API and check cache

#### func (*Request) ExecWithContext

```go
func (a *Request) ExecWithContext(ctx context.Context, body string) ([]byte, error)
```
ExecWithContext executes the API call and check cache with context

#### func (Request) GetID

```go
func (a Request) GetID(body string) string
```
GetID generates a hash ID for the request
