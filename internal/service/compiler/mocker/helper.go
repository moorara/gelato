package mocker

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/moorara/gelato/internal/service/compiler"
)

// normalizeMethods clones a field list of methods and substitude embedded interfaces with their methods.
func normalizeMethods(methods *ast.FieldList) *ast.FieldList {
	new := &ast.FieldList{}

	if methods == nil {
		return new
	}

	for _, method := range methods.List {
		// Embedded interface
		if len(method.Names) == 0 {
			// TODO:
		}

		// Method
		if len(method.Names) == 1 {
			if _, ok := method.Type.(*ast.FuncType); ok {
				new.List = append(new.List, &ast.Field{
					Names: method.Names,
					Type:  method.Type,
				})
			}
		}
	}

	return new
}

// normalizeFields clones a field list and converts embedded fields to non-embedded ones.
func normalizeFields(fields *ast.FieldList) *ast.FieldList {
	new := &ast.FieldList{}

	if fields == nil {
		return new
	}

	for _, field := range fields.List {
		f := &ast.Field{
			Names: field.Names,
			Type:  field.Type,
		}

		// Unnamed field
		if len(f.Names) == 0 {
			f.Names = []*ast.Ident{
				{Name: compiler.ConvertToUnexported(compiler.InferName(f.Type))},
			}
		}

		// Trailing arguments (for variadic functions)
		if e, ok := f.Type.(*ast.Ellipsis); ok {
			f.Type = &ast.ArrayType{
				Elt: e.Elt,
			}
		}

		new.List = append(new.List, f)
	}

	return new
}

// createKeyValueExprList creates a list of key-value assignments for creating structs from a field list.
func createKeyValueExprList(fieldList *ast.FieldList) []ast.Expr {
	list := []ast.Expr{}

	if fieldList == nil {
		return list
	}

	for _, f := range fieldList.List {
		if len(f.Names) > 0 {
			for _, n := range f.Names {
				list = append(list, &ast.KeyValueExpr{
					Key:   &ast.Ident{Name: n.Name},
					Value: &ast.Ident{Name: n.Name},
				})
			}
		} else {
			// Unnamed field
			name := compiler.ConvertToUnexported(compiler.InferName(f.Type))
			list = append(list, &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: name},
				Value: &ast.Ident{Name: name},
			})
		}
	}

	return list
}

func createZeroValueExpr(typ ast.Expr) ast.Expr {
	switch e := typ.(type) {
	case *ast.Ident:
		switch e.Name {
		case "error":
			return &ast.Ident{Name: "nil"}
		case "bool":
			return &ast.Ident{Name: "false"}
		case "string":
			return &ast.BasicLit{Kind: token.STRING, Value: `""`}
		case "byte", "rune":
			fallthrough
		case "int", "int8", "int16", "int32", "int64":
			fallthrough
		case "uint", "uint8", "uint16", "uint32", "uint64", "uintptr":
			return &ast.BasicLit{Kind: token.INT, Value: "0"}
		case "float32", "float64":
			return &ast.BasicLit{Kind: token.FLOAT, Value: "0.0"}
		case "complex64", "complex128":
			return &ast.BasicLit{Kind: token.IMAG, Value: "0.0i"}
		default: // struct
			return &ast.CompositeLit{Type: e}
		}

	case *ast.SelectorExpr:
		return &ast.CompositeLit{Type: e}

	case *ast.StarExpr, *ast.ArrayType, *ast.MapType, *ast.ChanType:
		return &ast.Ident{Name: "nil"}
	}

	panic(fmt.Sprintf("unknown type %T", typ))
}
