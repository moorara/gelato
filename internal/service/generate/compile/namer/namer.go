package namer

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

var (
	re1 = regexp.MustCompile(`^[a-z]`)
	re2 = regexp.MustCompile(`^[A-Z]+$`)
	re3 = regexp.MustCompile(`^[A-Z][0-9a-z_]`)
	re4 = regexp.MustCompile(`^([A-Z]+)[A-Z][0-9a-z_]`)
)

// IsExported determines whether or not a given name is exported.
func IsExported(name string) bool {
	first := name[0:1]
	return first == strings.ToUpper(first)
}

// ConvertToUnexported
func ConvertToUnexported(name string) string {
	switch {
	// Unexported (e.g. internal)
	case re1.MatchString(name):
		return name

	// All in upper letters (e.g. ID)
	case re2.MatchString(name):
		return strings.ToLower(name)

	// Starts with Title case (e.g. Request)
	case re3.MatchString(name):
		return strings.ToLower(name[0:1]) + name[1:]

	// Starts with all upper letters followed by a Title case (e.g. HTTPRequest)
	case re4.MatchString(name):
		m := re4.FindStringSubmatch(name)
		if len(m) == 2 {
			l := len(m[1])
			return strings.ToLower(name[0:l]) + name[l:]
		}
	}

	panic(fmt.Sprintf("ConvertToUnexported: unexpected identifer: %s", name))
}

// InferName
func InferName(expr ast.Expr) string {
	// First, get the last identifier name
	var lastName string
	astutil.Apply(expr,
		func(c *astutil.Cursor) bool {
			if id, ok := c.Node().(*ast.Ident); ok {
				lastName = id.Name
			}
			return true
		},
		func(c *astutil.Cursor) bool {
			return true
		},
	)

	return lastName
}
