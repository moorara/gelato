package decorate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	d := New()

	assert.NotNil(t, d)
}

func TestDecorator_Decorate(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		expectedError error
	}{
		{
			name:          "OK",
			path:          ".",
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d := &Decorator{}

			err := d.Decorate(tc.path)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
