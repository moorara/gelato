package modifier

import (
	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

// Modifier is an interface for a type that modifies an AST using the astutil.Apply function.
type Modifier interface {
	Pre(*astutil.Cursor) bool
	Post(*astutil.Cursor) bool
}

type modifier struct {
	depth  int
	logger *log.ColorfulLogger
}
