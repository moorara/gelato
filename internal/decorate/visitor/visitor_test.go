package visitor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVisitor(t *testing.T) {
	tests := []struct {
		name         string
		visitor      visitor
		expectedNext visitor
	}{
		{
			name: "OK",
			visitor: visitor{
				depth:  2,
				logger: nil,
			},
			expectedNext: visitor{
				depth:  3,
				logger: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			next := tc.visitor.Next()

			assert.Equal(t, tc.expectedNext, next)
		})
	}
}
