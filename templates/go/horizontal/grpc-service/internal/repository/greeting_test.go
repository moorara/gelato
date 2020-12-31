package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewGreetingRepository(t *testing.T) {
	tests := []struct {
		name          string
		expectedError error
	}{
		{
			name:          "OK",
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository, err := NewGreetingRepository()

			assert.NoError(t, err)
			assert.NotNil(t, repository)
		})
	}
}

func TestGreetingRepositoryString(t *testing.T) {
	tests := []struct {
		name           string
		repository     GreetingRepository
		expectedString string
	}{
		{
			name:           "OK",
			repository:     &greetingRepository{},
			expectedString: "greeting-repository",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.repository.String()

			assert.Equal(t, tc.expectedString, str)
		})
	}
}

func TestGreetingRepositoryConnect(t *testing.T) {
	tests := []struct {
		name          string
		repository    GreetingRepository
		expectedError string
	}{
		{
			name:          "Successful",
			repository:    &greetingRepository{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.repository.Connect()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGreetingRepositoryDisconnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	tests := []struct {
		name          string
		repository    GreetingRepository
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Successful",
			repository:    &greetingRepository{},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name:          "ContextExpires",
			repository:    &greetingRepository{},
			ctx:           ctx,
			expectedError: "context deadline exceeded",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.repository.Disconnect(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGreetingRepositoryCheckHealth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	tests := []struct {
		name          string
		repository    GreetingRepository
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Successful",
			repository:    &greetingRepository{},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name:          "ContextExpires",
			repository:    &greetingRepository{},
			ctx:           ctx,
			expectedError: "context deadline exceeded",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.repository.CheckHealth(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestGreetingRepositoryCreate(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		greeting      string
		expectedError error
	}{
		{
			name:          "EmptyGreeting",
			ctx:           context.Background(),
			greeting:      "",
			expectedError: errors.New("cannot store empty greeting"),
		},
		{
			name:          "OK",
			ctx:           context.Background(),
			greeting:      "Hello, World!",
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := &greetingRepository{
				store: map[time.Time]string{},
			}

			_, err := repository.Create(tc.ctx, tc.greeting)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGreetingRepositoryGet(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2020-01-08T06:14:57+04:30")
	t2 := time.Now()

	tests := []struct {
		name             string
		ctx              context.Context
		timestamp        time.Time
		expectedGreeting string
		expectedError    error
	}{
		{
			name:             "NotFound",
			ctx:              context.Background(),
			timestamp:        t1,
			expectedGreeting: "",
			expectedError:    errors.New("no greeting found for Wed, 08 Jan 2020 06:14:57 +0430"),
		},
		{
			name:             "Success",
			ctx:              context.Background(),
			timestamp:        t2,
			expectedGreeting: "Hello, World!",
			expectedError:    nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			repository := &greetingRepository{
				store: map[time.Time]string{
					t2: "Hello, World!",
				},
			}

			greeting, err := repository.Get(tc.ctx, tc.timestamp)

			assert.Equal(t, tc.expectedGreeting, greeting)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
