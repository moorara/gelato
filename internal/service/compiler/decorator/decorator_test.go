package decorator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/service/compiler"
	"github.com/moorara/gelato/internal/log"
)

func TestNew(t *testing.T) {
	c := New(log.Info)

	assert.NotNil(t, c)
	assert.IsType(t, &compiler.Compiler{}, c)
}
