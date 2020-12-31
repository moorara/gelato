package mapper

import (
	"errors"
	"testing"

	"horizontal/http-service/internal/entity"
	"horizontal/http-service/internal/idl"

	"github.com/stretchr/testify/assert"
)

func TestGreetRequestIDLToDomain(t *testing.T) {
	tests := []struct {
		name          string
		req           *idl.GreetRequest
		expectedReq   *entity.GreetRequest
		expectedError error
	}{
		{
			name:          "NilRequest",
			req:           nil,
			expectedReq:   nil,
			expectedError: errors.New("greet request cannot be nil"),
		},
		{
			name:          "EmptyName",
			req:           &idl.GreetRequest{},
			expectedReq:   nil,
			expectedError: errors.New("name cannot be empty"),
		},
		{
			name: "OK",
			req: &idl.GreetRequest{
				Name: "Jane",
			},
			expectedReq: &entity.GreetRequest{
				Name: "Jane",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req, err := GreetRequestIDLToDomain(tc.req)

			assert.Equal(t, tc.expectedReq, req)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGreetResponseDomainToIDL(t *testing.T) {
	tests := []struct {
		name          string
		resp          *entity.GreetResponse
		expectedResp  *idl.GreetResponse
		expectedError error
	}{
		{
			name:          "NilResponse",
			resp:          nil,
			expectedResp:  nil,
			expectedError: errors.New("greet response cannot be nil"),
		},
		{
			name:          "EmptyGreeting",
			resp:          &entity.GreetResponse{},
			expectedResp:  nil,
			expectedError: errors.New("greeting cannot be empty"),
		},
		{
			name: "OK",
			resp: &entity.GreetResponse{
				Greeting: "Hello, Jane!",
			},
			expectedResp: &idl.GreetResponse{
				Greeting: "Hello, Jane!",
			},
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := GreetResponseDomainToIDL(tc.resp)

			assert.Equal(t, tc.expectedResp, resp)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
