package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"horizontal/grpc-service/internal/handler"
)

type ServeMock struct {
	InLis    net.Listener
	OutError error
}

// MockGRPCServer is a mock implementation of grpcServer interface.
type MockGRPCServer struct {
	GracefulStopCounter int

	ServeCounter int
	ServeMocks   []ServeMock
}

func (m *MockGRPCServer) GracefulStop() {
	m.GracefulStopCounter++
}

func (m *MockGRPCServer) Serve(lis net.Listener) error {
	i := m.ServeCounter
	m.ServeCounter++
	m.ServeMocks[i].InLis = lis
	return m.ServeMocks[i].OutError
}

func TestNewGRPCServer(t *testing.T) {
	tests := []struct {
		name            string
		greetingHandler handler.GreetingHandler
		opts            GRPCServerOptions
	}{
		{
			name:            "OK",
			greetingHandler: &MockGreetingHandler{},
			opts:            GRPCServerOptions{},
		},
		{
			name:            "WithTLS",
			greetingHandler: &MockGreetingHandler{},
			opts: GRPCServerOptions{
				TLSCert:  &tls.Certificate{},
				ClientCA: x509.NewCertPool(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			server, err := NewGRPCServer(tc.greetingHandler, tc.opts)

			assert.NoError(t, err)
			assert.NotNil(t, server)
		})
	}
}

func TestGRPCServerString(t *testing.T) {
	tests := []struct {
		name           string
		server         *GRPCServer
		expectedString string
	}{
		{
			name:           "OK",
			server:         &GRPCServer{},
			expectedString: "grpc-server",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.server.String()

			assert.Equal(t, tc.expectedString, str)
		})
	}
}

func TestGRPCServerListenAndServe(t *testing.T) {
	tests := []struct {
		name          string
		server        *GRPCServer
		expectedError string
	}{
		{
			name: "ListenFails",
			server: &GRPCServer{
				addr: ":-1",
			},
			expectedError: "listen tcp: address -1: invalid port",
		},
		{
			name: "Successful",
			server: &GRPCServer{
				addr: "127.0.0.1:",
				server: &MockGRPCServer{
					ServeMocks: []ServeMock{
						{OutError: nil},
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

func TestGRPCServerShutdown(t *testing.T) {
	cancelledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name          string
		server        *GRPCServer
		ctx           context.Context
		expectedError string
	}{
		{
			name: "Successful",
			server: &GRPCServer{
				addr:   "127.0.0.1:",
				server: &MockGRPCServer{},
			},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name: "ContextCancelled",
			server: &GRPCServer{
				addr:   "127.0.0.1:",
				server: &MockGRPCServer{},
			},
			ctx:           cancelledCtx,
			expectedError: "context canceled",
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
