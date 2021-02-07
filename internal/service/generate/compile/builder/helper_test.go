package builder

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFieldInitExpr(t *testing.T) {
	tests := []struct {
		name         string
		id           *ast.Ident
		typ          ast.Expr
		expectedExpr *ast.KeyValueExpr
	}{
		{
			name: "error",
			id:   &ast.Ident{Name: "e"},
			typ:  &ast.Ident{Name: "error"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "e"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "bool",
			id:   &ast.Ident{Name: "b"},
			typ:  &ast.Ident{Name: "bool"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "b"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "string",
			id:   &ast.Ident{Name: "s"},
			typ:  &ast.Ident{Name: "string"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "s"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "byte",
			id:   &ast.Ident{Name: "b"},
			typ:  &ast.Ident{Name: "byte"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "b"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "rune",
			id:   &ast.Ident{Name: "r"},
			typ:  &ast.Ident{Name: "rune"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "r"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "int",
			id:   &ast.Ident{Name: "i"},
			typ:  &ast.Ident{Name: "int"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "int8",
			id:   &ast.Ident{Name: "i"},
			typ:  &ast.Ident{Name: "int8"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "int16",
			id:   &ast.Ident{Name: "i"},
			typ:  &ast.Ident{Name: "int16"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "int32",
			id:   &ast.Ident{Name: "i"},
			typ:  &ast.Ident{Name: "int32"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "int64",
			id:   &ast.Ident{Name: "i"},
			typ:  &ast.Ident{Name: "int64"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "i"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "uint",
			id:   &ast.Ident{Name: "u"},
			typ:  &ast.Ident{Name: "uint"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "u"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "uint8",
			id:   &ast.Ident{Name: "u"},
			typ:  &ast.Ident{Name: "uint8"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "u"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "uint16",
			id:   &ast.Ident{Name: "u"},
			typ:  &ast.Ident{Name: "uint16"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "u"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "uint32",
			id:   &ast.Ident{Name: "u"},
			typ:  &ast.Ident{Name: "uint32"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "u"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "uint64",
			id:   &ast.Ident{Name: "u"},
			typ:  &ast.Ident{Name: "uint64"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "u"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "uintptr",
			id:   &ast.Ident{Name: "u"},
			typ:  &ast.Ident{Name: "uintptr"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "u"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "float32",
			id:   &ast.Ident{Name: "f"},
			typ:  &ast.Ident{Name: "float32"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "f"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "float64",
			id:   &ast.Ident{Name: "f"},
			typ:  &ast.Ident{Name: "float64"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "f"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "complex64",
			id:   &ast.Ident{Name: "c"},
			typ:  &ast.Ident{Name: "complex64"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "c"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "complex128",
			id:   &ast.Ident{Name: "c"},
			typ:  &ast.Ident{Name: "complex128"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "c"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "Struct_SamePackage",
			id:   &ast.Ident{Name: "a"},
			typ:  &ast.Ident{Name: "Address"},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "a"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "Struct_OtherPackage",
			id:   &ast.Ident{Name: "t"},
			typ: &ast.SelectorExpr{
				X:   &ast.Ident{Name: "http"},
				Sel: &ast.Ident{Name: "Transport"},
			},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "t"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "Pointer",
			id:   &ast.Ident{Name: "p"},
			typ: &ast.StarExpr{
				X: &ast.Ident{Name: "int"},
			},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "p"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "Slice",
			id:   &ast.Ident{Name: "s"},
			typ: &ast.ArrayType{
				Elt: &ast.Ident{Name: "int"},
			},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "s"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "Map",
			id:   &ast.Ident{Name: "m"},
			typ: &ast.MapType{
				Key:   &ast.Ident{Name: "int"},
				Value: &ast.Ident{Name: "string"},
			},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "m"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
		{
			name: "Channel",
			id:   &ast.Ident{Name: "c"},
			typ: &ast.ChanType{
				Value: &ast.Ident{Name: "error"},
			},
			expectedExpr: &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: "c"},
				Value: &ast.Ident{Name: "nil"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			expr := createFieldInitExpr(tc.id, tc.typ)

			assert.Equal(t, tc.expectedExpr, expr)
		})
	}
}
