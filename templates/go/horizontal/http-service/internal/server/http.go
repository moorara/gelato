package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"horizontal/http-service/internal/handler"
	"horizontal/http-service/internal/idl"
	"horizontal/http-service/pkg/xhttp"
)

const (
	defaultHTTPPort = 4000
)

// httpServer is an interface for http.Server struct.
type httpServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HTTPServer is an HTTP server implementing graceful.Server interface.
type HTTPServer struct {
	server httpServer
}

// HTTPServerOptions are optional settings for creating an HTTP server.
type HTTPServerOptions struct {
	// The port number for the HTTP server.
	// The default port number is 8080.
	Port uint16

	// HTTP middleware for handlers.
	Middleware []xhttp.Middleware
}

// NewHTTPServer creates a new instance of HTTP Server.
func NewHTTPServer(healthHandler http.Handler, greetingHandler handler.GreetingHandler, opts HTTPServerOptions) (*HTTPServer, error) {
	if opts.Port == 0 {
		opts.Port = defaultHTTPPort
	}

	router := mux.NewRouter()
	router.Path("/health").Handler(healthHandler)
	idl.RegisterGreetingHandler(router, greetingHandler, opts.Middleware...)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", opts.Port),
		Handler: router,
	}

	return &HTTPServer{
		server: server,
	}, nil
}

// String returns the name of the server.
func (s *HTTPServer) String() string {
	return "http-server"
}

// ListenAndServe starts listening for incoming requests synchronously.
// It blocks the current goroutine until an error is returned.
func (s *HTTPServer) ListenAndServe() error {
	// Synchronous/Blocking
	// ListenAndServe always returns a non-nil error
	// After Shutdown or Close, the returned error is ErrServerClosed
	err := s.server.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server.
// It stops accepting new conenctions and blocks the current goroutine until all the pending requests are completed.
// If the context is cancelled, an error will be returned.
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
