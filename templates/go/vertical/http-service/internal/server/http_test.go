package server

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"vertical/http-service/internal/greeting"
)

type ListenAndServeMock struct {
	OutError error
}

type ShutdownMock struct {
	InCtx    context.Context
	OutError error
}

// MockHTTPServer is a mock implementation of httpServer interface.
type MockHTTPServer struct {
	ListenAndServeCounter int
	ListenAndServeMocks   []ListenAndServeMock

	ShutdownCounter int
	ShutdownMocks   []ShutdownMock
}

func (m *MockHTTPServer) ListenAndServe() error {
	i := m.ListenAndServeCounter
	m.ListenAndServeCounter++
	return m.ListenAndServeMocks[i].OutError
}

func (m *MockHTTPServer) Shutdown(ctx context.Context) error {
	i := m.ShutdownCounter
	m.ShutdownCounter++
	m.ShutdownMocks[i].InCtx = ctx
	return m.ShutdownMocks[i].OutError
}

func TestNewHTTPServer(t *testing.T) {
	tests := []struct {
		name            string
		healthHandler   http.Handler
		greetingService *greeting.Service
		opts            HTTPServerOptions
	}{
		{
			name: "OK",
			healthHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			}),
			greetingService: &greeting.Service{},
			opts:            HTTPServerOptions{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server, err := NewHTTPServer(tc.healthHandler, tc.greetingService, tc.opts)

			assert.NoError(t, err)
			assert.NotNil(t, server)
		})
	}
}

func TestHTTPServerString(t *testing.T) {
	tests := []struct {
		name           string
		server         *HTTPServer
		expectedString string
	}{
		{
			name:           "OK",
			server:         &HTTPServer{},
			expectedString: "http-server",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.server.String()

			assert.Equal(t, tc.expectedString, str)
		})
	}
}

func TestHTTPServerListenAndServe(t *testing.T) {
	tests := []struct {
		name          string
		server        *HTTPServer
		expectedError string
	}{
		{
			name: "ListenFails",
			server: &HTTPServer{
				server: &MockHTTPServer{
					ListenAndServeMocks: []ListenAndServeMock{
						{OutError: errors.New("error on listening")},
					},
				},
			},
			expectedError: "error on listening",
		},
		{
			name: "ServerClosed",
			server: &HTTPServer{
				server: &MockHTTPServer{
					ListenAndServeMocks: []ListenAndServeMock{
						{OutError: http.ErrServerClosed},
					},
				},
			},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.server.ListenAndServe()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestHTTPServerShutdown(t *testing.T) {
	tests := []struct {
		name          string
		server        *HTTPServer
		ctx           context.Context
		expectedError string
	}{
		{
			name: "Successful",
			server: &HTTPServer{
				server: &MockHTTPServer{
					ShutdownMocks: []ShutdownMock{
						{OutError: nil},
					},
				},
			},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name: "Unsuccessful",
			server: &HTTPServer{
				server: &MockHTTPServer{
					ShutdownMocks: []ShutdownMock{
						{OutError: errors.New("error on shutdown")},
					},
				},
			},
			ctx:           context.Background(),
			expectedError: "error on shutdown",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.server.Shutdown(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
