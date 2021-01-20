package example

// Foo is a struct.
type Foo struct {
	ID string
}

// Bar is an interface.
type Bar interface {
	Get(string) (string, error)
}
