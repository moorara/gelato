package mocker

import (
	"go/ast"

	"github.com/moorara/gelato/internal/service/generate/compile/namer"
)

func isEmbeddedInterface(method *ast.Field) bool {
	return len(method.Names) == 0
}

func isMethod(method *ast.Field) bool {
	if len(method.Names) == 1 {
		if _, ok := method.Type.(*ast.FuncType); ok {
			return true
		}
	}

	return false
}

// normalizeFieldList clones a field list and converts embedded fields to non-embedded ones.
func normalizeFieldList(fieldList *ast.FieldList) *ast.FieldList {
	new := &ast.FieldList{}

	for _, field := range fieldList.List {
		f := &ast.Field{
			Names: field.Names,
			Type:  field.Type,
		}

		// Unnamed field
		if len(f.Names) == 0 {
			f.Names = []*ast.Ident{
				{Name: namer.ConvertToUnexported(namer.InferName(f.Type))},
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
			name := namer.ConvertToUnexported(namer.InferName(f.Type))
			list = append(list, &ast.KeyValueExpr{
				Key:   &ast.Ident{Name: name},
				Value: &ast.Ident{Name: name},
			})
		}
	}

	return list
}
