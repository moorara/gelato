package gateway

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewTranslateGateway(t *testing.T) {
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
			gateway, err := NewTranslateGateway()

			assert.NoError(t, err)
			assert.NotNil(t, gateway)
		})
	}
}

func TestTranslateGatewayString(t *testing.T) {
	tests := []struct {
		name           string
		gateway        TranslateGateway
		expectedString string
	}{
		{
			name:           "OK",
			gateway:        &translateGateway{},
			expectedString: "translate-gateway",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			str := tc.gateway.String()

			assert.Equal(t, tc.expectedString, str)
		})
	}
}

func TestTranslateGatewayConnect(t *testing.T) {
	tests := []struct {
		name          string
		gateway       TranslateGateway
		expectedError string
	}{
		{
			name:          "Successful",
			gateway:       &translateGateway{},
			expectedError: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.gateway.Connect()

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTranslateGatewayDisconnect(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	tests := []struct {
		name          string
		gateway       TranslateGateway
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Successful",
			gateway:       &translateGateway{},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name:          "ContextExpires",
			gateway:       &translateGateway{},
			ctx:           ctx,
			expectedError: "context deadline exceeded",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.gateway.Disconnect(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTranslateGatewayCheckHealth(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	tests := []struct {
		name          string
		gateway       TranslateGateway
		ctx           context.Context
		expectedError string
	}{
		{
			name:          "Successful",
			gateway:       &translateGateway{},
			ctx:           context.Background(),
			expectedError: "",
		},
		{
			name:          "ContextExpires",
			gateway:       &translateGateway{},
			ctx:           ctx,
			expectedError: "context deadline exceeded",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.gateway.CheckHealth(tc.ctx)

			if tc.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError)
			}
		})
	}
}

func TestTranslateGatewayGetValue(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		phraseID       string
		languageCode   string
		expectedString string
		expectedError  error
	}{
		{
			name:           "EmptyPhraseID",
			ctx:            context.Background(),
			phraseID:       "",
			languageCode:   "",
			expectedString: "",
			expectedError:  errors.New("invalid phrase id"),
		},
		{
			name:           "EmptyLanguageCode",
			ctx:            context.Background(),
			phraseID:       "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			languageCode:   "",
			expectedString: "",
			expectedError:  errors.New("invalid language code"),
		},
		{
			name:           "OK",
			ctx:            context.Background(),
			phraseID:       "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			languageCode:   "en",
			expectedString: "Hello",
			expectedError:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gateway := &translateGateway{}
			str, err := gateway.GetValue(tc.ctx, tc.phraseID, tc.languageCode)

			assert.Equal(t, tc.expectedString, str)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
