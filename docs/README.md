# Documentation

## Test API

### Example

```go
package lookuptest

type Request struct {
  ID string
}

type Response struct {
  Name string
}

type Service interface {
  Lookup(*Request) (*Response, error)
}
```

```go
package handler

import (
  "io"
  "net/http"

  "github.com/octocat/example/services/lookup"
)

type Handler struct {
  lookupService lookup.Service
}

func New(ls lookup.Service) *Handler {
  return &Handler{
    lookupService: ls,
  }
}

func (h *Handler) Greet(w http.ResponseWriter, r *http.Request) {
  greeting, err := h.lookupService.Lookup()
  if err != nil {
    return
  }
  io.WriteString(w, greeting)
}
```

```go
package handler

import (
  "time"
  "testing"

  "github.com/octocat/example/.gen/services/lookuptest"
)

func TestHandler_Greet(t *testing.T) {
  req := lookuptest.Request()
  req = lookuptest.BuildRequest().WithID("1234").Value()
  req = lookuptest.BuildRequest().WithID("1234").Pointer()

  resp := lookuptest.Response()
  resp = lookuptest.BuildResponse().WithName("Jane Doe").Value()
  resp = lookuptest.BuildResponse().WithName("Jane Doe").Pointer()

  mock := lookuptest.MockService()
  mock.Expect().Lookup().WithArgs(req).Return(resp, nil)
  mock.Expect().Lookup().WithArgs(req).Call(func(req *LookupRequest) (*LookupResponse, error){
    time.Sleep(time.Second)
    return rep, nil
  })

  ls := mock.Impl()
  h := New(ls)

  w := httptest.NewRecorder()
  r := httptest.NewRequest("GET", "http://example.com", nil)
  h.Greet(w, r)

  mock.Assert(t)
}
```
