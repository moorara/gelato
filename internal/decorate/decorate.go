package decorate

import (
	"github.com/moorara/color"
)

// Decorator decorates a Go application.
type Decorator struct {
}

// New creates a new decorator.
func New() *Decorator {
	return &Decorator{}
}

// Decorate decorates a Go application.
func (d *Decorator) Decorate(path string) error {
	color.White("Decorating ...")
	return nil
}
