// Package client defines the interface for a generic client (external service, database, message queue, etc.).
// It also provides a mock implementation of the Client interface.
package client

import (
	"context"
	"time"

	"github.com/moorara/graceful"
	"github.com/moorara/health"
)

// Client is the interface for a generic client (external service, database, message queue, etc.).
type Client interface {
	graceful.Client
	health.Checker
}

// client implements the Client interface.
type client struct {
	name string
}

// New creates a new client.
func New(name string) (Client, error) {
	return &client{
		name: name,
	}, nil
}

// String returns the name of the client.
func (c *client) String() string {
	return c.name
}

// Connect connects the client by opening a long-lived connection.
func (c *client) Connect() error {
	time.Sleep(time.Second)
	return nil
}

// Disconnect disconnects the client by closing the long-lived connection.
// If the context is cancelled, an error will be returned.
func (c *client) Disconnect(ctx context.Context) error {
	select {
	case <-time.After(time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// CheckHealth checks the health of connection.
// If the context is cancelled, an error will be returned.
func (c *client) CheckHealth(ctx context.Context) error {
	select {
	case <-time.After(50 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
