package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreetRequest(t *testing.T) {
	tests := []struct {
		name           string
		entity         GreetRequest
		expectedString string
	}{
		{
			name: "OK",
			entity: GreetRequest{
				Name: "Jane",
			},
			expectedString: "GreetRequest{name=Jane}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.entity.String())
		})
	}
}

func TestGreetResponse(t *testing.T) {
	tests := []struct {
		name           string
		entity         GreetResponse
		expectedString string
	}{
		{
			name: "OK",
			entity: GreetResponse{
				Greeting: "Hello, Jane!",
			},
			expectedString: "GreetResponse{greeting=Hello, Jane!}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.entity.String())
		})
	}
}
