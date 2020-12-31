package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		clientName    string
		expectedError error
	}{
		{
			name:          "OK",
			clientName:    "test-client",
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client, err := New(tc.clientName)

			assert.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}

func TestClientString(t *testing.T) {
	tests := []struct {
		name           string
		client         Client
		expectedString string
	}{
		{
			name: "OK",
			client: &client{
				name: "test-client",
			},
			expectedString: "test-client",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.client.String()

			assert.Equal(t, tc.expectedString, str)
		})
	}
}

func TestClientConnect(t *testing.T) {
	tests := []struct {
		name          string
		client        Client
		expectedError string
	}{
		{
			name:          "Successful",
			client:        &client{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.client.Connect()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestClientDisconnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	tests := []struct {
		name          string
		client        Client
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Successful",
			client:        &client{},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name:          "ContextExpires",
			client:        &client{},
			ctx:           ctx,
			expectedError: "context deadline exceeded",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.client.Disconnect(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestClientCheckHealth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	tests := []struct {
		name          string
		client        Client
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Successful",
			client:        &client{},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name:          "ContextExpires",
			client:        &client{},
			ctx:           ctx,
			expectedError: "context deadline exceeded",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.client.CheckHealth(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}
