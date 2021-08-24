package builder

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateFieldInitDefs(t *testing.T) {
	tests := []struct {
		name          string
		id            *ast.Ident
		typ           ast.Expr
		expectedStmts []ast.Stmt
		expectedExpr  ast.Expr
	}{
		{
			name:          "string",
			id:            &ast.Ident{Name: "s"},
			typ:           &ast.Ident{Name: "string"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "String"},
				},
			},
		},
		{
			name:          "bool",
			id:            &ast.Ident{Name: "b"},
			typ:           &ast.Ident{Name: "bool"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Bool"},
				},
			},
		},
		{
			name:          "byte",
			id:            &ast.Ident{Name: "b"},
			typ:           &ast.Ident{Name: "byte"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Byte"},
				},
			},
		},
		{
			name:          "rune",
			id:            &ast.Ident{Name: "r"},
			typ:           &ast.Ident{Name: "rune"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Rune"},
				},
			},
		},
		{
			name:          "int",
			id:            &ast.Ident{Name: "i"},
			typ:           &ast.Ident{Name: "int"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int"},
				},
			},
		},
		{
			name:          "int8",
			id:            &ast.Ident{Name: "i"},
			typ:           &ast.Ident{Name: "int8"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int8"},
				},
			},
		},
		{
			name:          "int16",
			id:            &ast.Ident{Name: "i"},
			typ:           &ast.Ident{Name: "int16"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int16"},
				},
			},
		},
		{
			name:          "int32",
			id:            &ast.Ident{Name: "i"},
			typ:           &ast.Ident{Name: "int32"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int32"},
				},
			},
		},
		{
			name:          "int64",
			id:            &ast.Ident{Name: "i"},
			typ:           &ast.Ident{Name: "int64"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int64"},
				},
			},
		},
		{
			name:          "uint",
			id:            &ast.Ident{Name: "u"},
			typ:           &ast.Ident{Name: "uint"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint"},
				},
			},
		},
		{
			name:          "uint8",
			id:            &ast.Ident{Name: "u"},
			typ:           &ast.Ident{Name: "uint8"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint8"},
				},
			},
		},
		{
			name:          "uint16",
			id:            &ast.Ident{Name: "u"},
			typ:           &ast.Ident{Name: "uint16"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint16"},
				},
			},
		},
		{
			name:          "uint32",
			id:            &ast.Ident{Name: "u"},
			typ:           &ast.Ident{Name: "uint32"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint32"},
				},
			},
		},
		{
			name:          "uint64",
			id:            &ast.Ident{Name: "u"},
			typ:           &ast.Ident{Name: "uint64"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint64"},
				},
			},
		},
		{
			name:          "uintptr",
			id:            &ast.Ident{Name: "u"},
			typ:           &ast.Ident{Name: "uintptr"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uintptr"},
				},
			},
		},
		{
			name:          "float32",
			id:            &ast.Ident{Name: "f"},
			typ:           &ast.Ident{Name: "float32"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Float32"},
				},
			},
		},
		{
			name:          "float64",
			id:            &ast.Ident{Name: "f"},
			typ:           &ast.Ident{Name: "float64"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Float64"},
				},
			},
		},
		{
			name:          "complex64",
			id:            &ast.Ident{Name: "c"},
			typ:           &ast.Ident{Name: "complex64"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Complex64"},
				},
			},
		},
		{
			name:          "complex128",
			id:            &ast.Ident{Name: "c"},
			typ:           &ast.Ident{Name: "complex128"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Complex128"},
				},
			},
		},
		{
			name:          "error",
			id:            &ast.Ident{Name: "e"},
			typ:           &ast.Ident{Name: "error"},
			expectedStmts: nil,
			expectedExpr: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Error"},
				},
			},
		},
		{
			name: "Pointer",
			id:   &ast.Ident{Name: "no"},
			typ: &ast.StarExpr{
				X: &ast.Ident{Name: "int"},
			},
			expectedStmts: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{Name: "noPtr"},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "factory"},
								Sel: &ast.Ident{Name: "Int"},
							},
						},
					},
				},
			},
			expectedExpr: &ast.UnaryExpr{
				Op: token.AND,
				X:  &ast.Ident{Name: "noPtr"},
			},
		},
		{
			name: "Slice",
			id:   &ast.Ident{Name: "values"},
			typ: &ast.ArrayType{
				Elt: &ast.Ident{Name: "string"},
			},
			expectedStmts: []ast.Stmt{},
			expectedExpr: &ast.CompositeLit{
				Type: &ast.ArrayType{
					Elt: &ast.Ident{Name: "string"},
				},
				Elts: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "factory"},
							Sel: &ast.Ident{Name: "String"},
						},
					},
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "factory"},
							Sel: &ast.Ident{Name: "String"},
						},
					},
				},
			},
		},
		{
			name: "Map",
			id:   &ast.Ident{Name: "cache"},
			typ: &ast.MapType{
				Key:   &ast.Ident{Name: "int"},
				Value: &ast.Ident{Name: "string"},
			},
			expectedStmts: []ast.Stmt{},
			expectedExpr: &ast.CompositeLit{
				Type: &ast.MapType{
					Key:   &ast.Ident{Name: "int"},
					Value: &ast.Ident{Name: "string"},
				},
				Elts: []ast.Expr{
					&ast.KeyValueExpr{
						Key: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "factory"},
								Sel: &ast.Ident{Name: "Int"},
							},
						},
						Value: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "factory"},
								Sel: &ast.Ident{Name: "String"},
							},
						},
					},
					&ast.KeyValueExpr{
						Key: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "factory"},
								Sel: &ast.Ident{Name: "Int"},
							},
						},
						Value: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "factory"},
								Sel: &ast.Ident{Name: "String"},
							},
						},
					},
				},
			},
		},
		{
			name: "Channel",
			id:   &ast.Ident{Name: "messages"},
			typ: &ast.ChanType{
				Value: &ast.Ident{Name: "string"},
			},
			expectedStmts: []ast.Stmt{
				// messages := make(chan string, 2)
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{Name: "messages"},
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.Ident{Name: "make"},
							Args: []ast.Expr{
								&ast.ChanType{
									Dir:   ast.SEND | ast.RECV,
									Value: &ast.Ident{Name: "string"},
								},
								&ast.BasicLit{
									Kind:  token.INT,
									Value: "2",
								},
							},
						},
					},
				},
				// messages <- factory.String()
				&ast.SendStmt{
					Chan: &ast.Ident{Name: "messages"},
					Value: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "factory"},
							Sel: &ast.Ident{Name: "String"},
						},
					},
				},
				// messages <- factory.String()
				&ast.SendStmt{
					Chan: &ast.Ident{Name: "messages"},
					Value: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "factory"},
							Sel: &ast.Ident{Name: "String"},
						},
					},
				},
			},
			expectedExpr: &ast.Ident{Name: "messages"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmts, expr := createFieldInitDefs(tc.id, tc.typ)

			assert.Equal(t, tc.expectedStmts, stmts)
			assert.Equal(t, tc.expectedExpr, expr)
		})
	}
}
