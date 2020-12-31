package greeting

import "net/http"

// MockMiddleware is a mock implementation for xhttp.Middleware interface.
type MockMiddleware struct{}

func (m *MockMiddleware) Wrap(next http.HandlerFunc) http.HandlerFunc {
	return next
}
