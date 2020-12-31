package gateway

import (
	"context"
	"errors"
	"time"

	"github.com/moorara/graceful"
	"github.com/moorara/health"
)

// TranslateGateway is the interface for calling an imaginary translation service.
type TranslateGateway interface {
	graceful.Client
	health.Checker
	GetValue(ctx context.Context, phraseID, languageCode string) (string, error)
}

// translateGateway implements TranslateGateway interface.
type translateGateway struct{}

// NewTranslateGateway creates a new instance of TranslateGateway.
func NewTranslateGateway() (TranslateGateway, error) {
	return &translateGateway{}, nil
}

// String returns the name of the gateway.
func (g *translateGateway) String() string {
	return "translate-gateway"
}

// Connect opens a long-lived connection to the external service.
func (g *translateGateway) Connect() error {
	time.Sleep(time.Second)
	return nil
}

// Disconnect closes the long-lived connection to the external service.
// If the context is cancelled, an error will be returned.
func (g *translateGateway) Disconnect(ctx context.Context) error {
	select {
	case <-time.After(time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// CheckHealth checks the health of connection to the external service.
// If the context is cancelled, an error will be returned.
func (g *translateGateway) CheckHealth(ctx context.Context) error {
	select {
	case <-time.After(50 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// GetValue returns the string for a phrase in a given language.
func (g *translateGateway) GetValue(ctx context.Context, phraseID, languageCode string) (string, error) {
	if phraseID == "" {
		return "", errors.New("invalid phrase id")
	}

	if languageCode == "" {
		return "", errors.New("invalid language code")
	}

	// Make a call to the service using the connection
	time.Sleep(100 * time.Millisecond)
	return "Hello", nil
}
