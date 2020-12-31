package server

import "net/http"

type GreetMock struct {
	InResponseWriter http.ResponseWriter
	InRequest        *http.Request
}

// MockGreetingHandler is a mock implementation for handler.GreetingHandler interface.
type MockGreetingHandler struct {
	GreetCounter int
	GreetMocks   []GreetMock
}

func (m *MockGreetingHandler) Greet(w http.ResponseWriter, r *http.Request) {
	i := m.GreetCounter
	m.GreetCounter++
	m.GreetMocks[i].InResponseWriter = w
	m.GreetMocks[i].InRequest = r
}
