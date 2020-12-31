package mapper

import (
	"errors"
	"testing"

	"horizontal/grpc-service/internal/entity"
	"horizontal/grpc-service/internal/idl/greetingpb"

	"github.com/stretchr/testify/assert"
)

func TestGreetRequestIDLToDomain(t *testing.T) {
	tests := []struct {
		name          string
		req           *greetingpb.GreetRequest
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
			req:           &greetingpb.GreetRequest{},
			expectedReq:   nil,
			expectedError: errors.New("name cannot be empty"),
		},
		{
			name: "OK",
			req: &greetingpb.GreetRequest{
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
		expectedResp  *greetingpb.GreetResponse
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
			expectedResp: &greetingpb.GreetResponse{
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
