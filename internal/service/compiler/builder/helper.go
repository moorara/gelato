package builder

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"github.com/moorara/gelato/internal/service/compiler"
)

const size = 2

func createFieldInitDefs(id *ast.Ident, typ ast.Expr) ([]ast.Stmt, ast.Expr) {
	unexportedName := compiler.ConvertToUnexported(id.Name)

	switch e := typ.(type) {
	case *ast.StarExpr:
		ptrID := &ast.Ident{Name: fmt.Sprintf("%sPtr", unexportedName)}
		stmts, expr := createFieldInitDefs(ptrID, e.X)
		stmts = append(stmts, &ast.AssignStmt{
			Lhs: []ast.Expr{ptrID},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{expr},
		})

		return stmts, &ast.UnaryExpr{
			Op: token.AND,
			X:  ptrID,
		}

	case *ast.ArrayType:
		stmts := []ast.Stmt{}
		elts := []ast.Expr{}

		for i := 0; i < size; i++ {
			s, expr := createFieldInitDefs(id, e.Elt)
			stmts = append(stmts, s...)
			elts = append(elts, expr)
		}

		return stmts, &ast.CompositeLit{
			Type: e,
			Elts: elts,
		}

	case *ast.MapType:
		stmts := []ast.Stmt{}
		elts := []ast.Expr{}

		for i := 0; i < size; i++ {
			keyStmts, keyExpr := createFieldInitDefs(id, e.Key)
			valStmts, valExpr := createFieldInitDefs(id, e.Value)
			stmts = append(stmts, keyStmts...)
			stmts = append(stmts, valStmts...)
			elts = append(elts, &ast.KeyValueExpr{
				Key:   keyExpr,
				Value: valExpr,
			})
		}

		return stmts, &ast.CompositeLit{
			Type: e,
			Elts: elts,
		}

	case *ast.ChanType:
		chanID := &ast.Ident{Name: unexportedName}
		stmts := []ast.Stmt{
			// id := make(chan ..., 2)
			&ast.AssignStmt{
				Lhs: []ast.Expr{chanID},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.Ident{Name: "make"},
						Args: []ast.Expr{
							&ast.ChanType{
								Dir:   ast.SEND | ast.RECV, // We always create a bi-directional channel, so we can always write to it and read from it
								Value: e.Value,
							},
							&ast.BasicLit{
								Kind:  token.INT,
								Value: strconv.Itoa(size),
							},
						},
					},
				},
			},
		}

		for i := 0; i < size; i++ {
			valStmts, valExpr := createFieldInitDefs(id, e.Value)
			stmts = append(stmts, valStmts...)
			stmts = append(stmts, &ast.SendStmt{
				Chan:  chanID,
				Value: valExpr,
			})
		}

		return stmts, chanID

	// Type in another package
	case *ast.SelectorExpr:
		// TODO:
		return nil, &ast.Ident{Name: "nil"}

	case *ast.Ident:
		switch e.Name {
		case "string":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "String"},
				},
			}

		case "bool":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Bool"},
				},
			}

		case "byte":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Byte"},
				},
			}

		case "rune":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Rune"},
				},
			}

		case "int":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int"},
				},
			}

		case "int8":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int8"},
				},
			}

		case "int16":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int16"},
				},
			}

		case "int32":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int32"},
				},
			}

		case "int64":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Int64"},
				},
			}

		case "uint":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint"},
				},
			}

		case "uint8":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint8"},
				},
			}

		case "uint16":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint16"},
				},
			}

		case "uint32":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint32"},
				},
			}

		case "uint64":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uint64"},
				},
			}

		case "uintptr":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Uintptr"},
				},
			}

		case "float32":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Float32"},
				},
			}

		case "float64":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Float64"},
				},
			}

		case "complex64":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Complex64"},
				},
			}

		case "complex128":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Complex128"},
				},
			}

		case "error":
			return nil, &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "factory"},
					Sel: &ast.Ident{Name: "Error"},
				},
			}

		// Type in the same package
		default:
			// TODO:
			return nil, &ast.Ident{Name: "nil"}
		}
	}

	panic(fmt.Sprintf("unknown type %T", typ))
}
