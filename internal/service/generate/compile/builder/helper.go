package builder

import (
	"go/ast"
)

func createFieldInitExpr(id *ast.Ident, typ ast.Expr) *ast.KeyValueExpr {
	var value ast.Expr

	// TODO:
	value = &ast.Ident{Name: "nil"}

	switch e := typ.(type) {
	case *ast.Ident:
		switch e.Name {
		case "error":
		case "bool":
		case "string":
		case "byte":
		case "rune":
		case "int":
		case "int8":
		case "int16":
		case "int32":
		case "int64":
		case "uint":
		case "uint8":
		case "uint16":
		case "uint32":
		case "uint64":
		case "uintptr":
		case "float32":
		case "float64":
		case "complex64":
		case "complex128":
		default: // struct
		}
	case *ast.SelectorExpr:
	case *ast.StarExpr:
	case *ast.ArrayType:
	case *ast.MapType:
	case *ast.ChanType:
	}

	return &ast.KeyValueExpr{
		Key:   &ast.Ident{Name: id.Name},
		Value: value,
	}
}
