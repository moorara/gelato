package mocker

import (
	"go/ast"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeFieldList(t *testing.T) {
	tests := []struct {
		name              string
		fieldList         *ast.FieldList
		expectedFieldList *ast.FieldList
	}{
		{
			name: "NamedFields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "s"},
						},
						Type: &ast.Ident{Name: "string"},
					},
				},
			},
			expectedFieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "s"},
						},
						Type: &ast.Ident{Name: "string"},
					},
				},
			},
		},
		{
			name: "UnnamedFields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{Name: "string"},
					},
				},
			},
			expectedFieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "string"},
						},
						Type: &ast.Ident{Name: "string"},
					},
				},
			},
		},
		{
			name: "TrailingFields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "strings"},
						},
						Type: &ast.Ellipsis{
							Elt: &ast.Ident{Name: "string"},
						},
					},
				},
			},
			expectedFieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "strings"},
						},
						Type: &ast.ArrayType{
							Elt: &ast.Ident{Name: "string"},
						},
					},
				},
			},
		},
		{
			name: "UnnamedTrailingFields",
			fieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ellipsis{
							Elt: &ast.Ident{Name: "string"},
						},
					},
				},
			},
			expectedFieldList: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "string"},
						},
						Type: &ast.ArrayType{
							Elt: &ast.Ident{Name: "string"},
						},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fieldList := normalizeFieldList(tc.fieldList)

			assert.Equal(t, tc.expectedFieldList, fieldList)
		})
	}
}
