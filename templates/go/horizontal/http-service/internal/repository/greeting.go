package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/moorara/graceful"
	"github.com/moorara/health"
)

// GreetingRepository is the interface for interacting with the data store for greetigns.
type GreetingRepository interface {
	graceful.Client
	health.Checker
	Create(ctx context.Context, greeting string) (time.Time, error)
	Get(ctx context.Context, t time.Time) (string, error)
}

// greetingRepository is an in-memory key-value data store for greetings that implements GreetingRepository interface.
type greetingRepository struct {
	store map[time.Time]string
}

// NewGreetingRepository creates a new repository for storing and retrieving greetings.
func NewGreetingRepository() (GreetingRepository, error) {
	return &greetingRepository{
		store: make(map[time.Time]string),
	}, nil
}

// String returns the name of the repository.
func (r *greetingRepository) String() string {
	return "greeting-repository"
}

// Connect opens a long-lived connection to the repository backend.
func (r *greetingRepository) Connect() error {
	time.Sleep(time.Second)
	return nil
}

// Disconnect closes the long-lived connection to the repository backend.
// If the context is cancelled, an error will be returned.
func (r *greetingRepository) Disconnect(ctx context.Context) error {
	select {
	case <-time.After(time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// CheckHealth checks the health of connection to the repository backend.
// If the context is cancelled, an error will be returned.
func (r *greetingRepository) CheckHealth(ctx context.Context) error {
	select {
	case <-time.After(50 * time.Millisecond):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Create stores a greeting and returns the creation timestamp.
func (r *greetingRepository) Create(ctx context.Context, greeting string) (time.Time, error) {
	if greeting == "" {
		return time.Time{}, errors.New("cannot store empty greeting")
	}

	t := time.Now()
	r.store[t] = greeting
	return t, nil
}

// Get retrieves a greeting by its creation timestamp.
func (r *greetingRepository) Get(ctx context.Context, t time.Time) (string, error) {
	greeting, ok := r.store[t]
	if !ok {
		return "", fmt.Errorf("no greeting found for %s", t.Format(time.RFC1123Z))
	}

	return greeting, nil
}
