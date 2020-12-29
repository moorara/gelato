package modifier

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/moorara/gelato/internal/log"
)

func TestGenericModifier(t *testing.T) {
	logger := log.New(log.None)
	clogger := &log.ColorfulLogger{
		Red:     logger,
		Green:   logger,
		Yellow:  logger,
		Blue:    logger,
		Magenta: logger,
		Cyan:    logger,
		White:   logger,
	}

	fileNode := &ast.File{
		Name: &ast.Ident{
			Name: "controller",
		},
		Decls: []ast.Decl{
			// Imports
			&ast.GenDecl{
				Tok: token.IMPORT,
				Specs: []ast.Spec{
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "context",
						},
					},
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "github.com/octocat/service/internal/entity",
						},
					},
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "github.com/octocat/service/internal/gateway",
						},
					},
					&ast.ImportSpec{
						Name: nil,
						Path: &ast.BasicLit{
							Value: "github.com/octocat/service/internal/repository",
						},
					},
				},
			},
			// Interface
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "Controller",
						},
						Type: &ast.InterfaceType{
							Methods: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{},
										Type:  &ast.FuncType{},
									},
								},
							},
						},
					},
				},
			},
			// Struct
			&ast.GenDecl{
				Tok: token.TYPE,
				Specs: []ast.Spec{
					&ast.TypeSpec{
						Name: &ast.Ident{
							Name: "controller",
						},
						Type: &ast.StructType{
							Fields: &ast.FieldList{
								List: []*ast.Field{
									{
										Names: []*ast.Ident{},
										Type:  &ast.SelectorExpr{},
									},
								},
							},
						},
					},
				},
			},
			// Exported Function
			&ast.FuncDecl{
				Recv: nil,
				Name: &ast.Ident{Name: "NewController"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "userGateway"},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "gateway"},
										Sel: &ast.Ident{Name: "UserGateway"},
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.Ident{Name: "Controller"},
							},
							{
								Type: &ast.Ident{Name: "error"},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.UnaryExpr{
									X: &ast.CompositeLit{
										Type: &ast.Ident{Name: "controller"},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key:   &ast.Ident{Name: "userGateway"},
												Value: &ast.Ident{Name: "userGateway"},
											},
										},
									},
								},
								&ast.Ident{Name: "nil"},
							},
						},
					},
				},
			},
			// Unexported Function
			&ast.FuncDecl{
				Recv: nil,
				Name: &ast.Ident{Name: "newController"},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "ug"},
								},
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "gateway"},
										Sel: &ast.Ident{Name: "UserGateway"},
									},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "controller"},
								},
							},
							{
								Type: &ast.Ident{Name: "error"},
							},
						},
					},
				},
			},
			// Exported Method
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{Name: "c"},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{Name: "controller"},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "Calculate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "a"},
									{Name: "b"},
								},
								Type: &ast.StarExpr{
									X: &ast.Ident{Name: "int"},
								},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "entity"},
										Sel: &ast.Ident{Name: "CalculateResponse"},
									},
								},
							},
							{
								Type: &ast.Ident{Name: "error"},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.UnaryExpr{
									X: &ast.CompositeLit{
										Type: &ast.StarExpr{
											X: &ast.SelectorExpr{
												X:   &ast.Ident{Name: "entity"},
												Sel: &ast.Ident{Name: "CalculateResponse"},
											},
										},
										Elts: []ast.Expr{
											&ast.KeyValueExpr{
												Key: &ast.Ident{Name: "Result"},
												Value: &ast.BinaryExpr{
													Op: token.Lookup("*"),
													X:  &ast.Ident{Name: "a"},
													Y:  &ast.Ident{Name: "b"},
												},
											},
										},
									},
								},
								&ast.Ident{Name: "nil"},
							},
						},
					},
				},
			},
			// UnexportedMethod
			&ast.FuncDecl{
				Recv: &ast.FieldList{
					List: []*ast.Field{
						{
							Names: []*ast.Ident{
								{Name: "c"},
							},
							Type: &ast.StarExpr{
								X: &ast.Ident{Name: "controller"},
							},
						},
					},
				},
				Name: &ast.Ident{
					Name: "calculate",
				},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "a"},
									{Name: "b"},
								},
								Type: &ast.Ident{Name: "int"},
							},
						},
					},
					Results: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "resp"},
								},
								Type: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "entity"},
									Sel: &ast.Ident{Name: "CalculateResponse"},
								},
							},
							{
								Names: []*ast.Ident{
									{Name: "err"},
								},
								Type: &ast.Ident{Name: "error"},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name         string
		depth        int
		module       string
		relPath      string
		node         ast.Node
		expectedNode ast.Node
	}{
		{
			name:         "OK",
			depth:        2,
			module:       "github.com/octocat/service",
			relPath:      "internal/controller",
			node:         fileNode,
			expectedNode: fileNode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := NewGeneric(tc.depth, clogger)

			node := m.Modify(tc.module, tc.relPath, tc.node)

			assert.Equal(t, tc.expectedNode, node)
		})
	}
}

func TestGenericImportModifier(t *testing.T) {
	logger := log.New(log.None)
	clogger := &log.ColorfulLogger{
		Red:     logger,
		Green:   logger,
		Yellow:  logger,
		Blue:    logger,
		Magenta: logger,
		Cyan:    logger,
		White:   logger,
	}

	importGenDecl := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{
				Name: &ast.Ident{
					Name: "u",
				},
				Path: &ast.BasicLit{
					Value: "net/url",
				},
			},
		},
	}

	tests := []struct {
		name         string
		depth        int
		origPkgName  string
		origPkgPath  string
		node         ast.Node
		expectedNode ast.Node
	}{
		{
			name:        "InvalidGenDecl",
			depth:       2,
			origPkgName: "_controller",
			origPkgPath: "github.com/octocat/service/internal/controller",
			node: &ast.GenDecl{
				Tok: token.TYPE,
			},
			expectedNode: &ast.GenDecl{
				Tok: token.TYPE,
			},
		},
		{
			name:         "ImportGenDecl",
			depth:        2,
			origPkgName:  "_controller",
			origPkgPath:  "github.com/octocat/service/internal/controller",
			node:         importGenDecl,
			expectedNode: importGenDecl,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &genericImportModifier{
				modifier: modifier{
					depth:  tc.depth,
					logger: clogger,
				},
			}

			node := m.Modify(tc.origPkgName, tc.origPkgPath, tc.node)

			assert.Equal(t, tc.expectedNode, node)
		})
	}
}

func TestGenericTypeModifier(t *testing.T) {
	logger := log.New(log.None)
	clogger := &log.ColorfulLogger{
		Red:     logger,
		Green:   logger,
		Yellow:  logger,
		Blue:    logger,
		Magenta: logger,
		Cyan:    logger,
		White:   logger,
	}

	interfaceGenDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: "Controller"},
				Type: &ast.InterfaceType{
					Methods: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "Calculate"},
								},
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.SelectorExpr{
													X:   &ast.Ident{Name: "context"},
													Sel: &ast.Ident{Name: "Context"},
												},
											},
											{
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "entity"},
														Sel: &ast.Ident{Name: "CalculateRequest"},
													},
												},
											},
										},
									},
									Results: &ast.FieldList{
										List: []*ast.Field{
											{
												Type: &ast.StarExpr{
													X: &ast.SelectorExpr{
														X:   &ast.Ident{Name: "entity"},
														Sel: &ast.Ident{Name: "CalculateResponse"},
													},
												},
											},
											{
												Type: &ast.Ident{Name: "error"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	structGenDecl := &ast.GenDecl{
		Tok: token.TYPE,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: "controller"},
				Type: &ast.StructType{
					Fields: &ast.FieldList{
						List: []*ast.Field{
							{
								Names: []*ast.Ident{
									{Name: "userGateway"},
								},
								Type: &ast.SelectorExpr{
									X:   &ast.Ident{Name: "gateway"},
									Sel: &ast.Ident{Name: "UserGateway"},
								},
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name              string
		depth             int
		origPkgName       string
		interfaceName     string
		node              ast.Node
		expectedNode      ast.Node
		expectedInterface *interfaceType
		expectedStruct    *structType
	}{
		{
			name:          "InvalidGenDecl",
			depth:         2,
			origPkgName:   "controller",
			interfaceName: "Controller",
			node: &ast.GenDecl{
				Tok: token.IMPORT,
			},
			expectedNode: &ast.GenDecl{
				Tok: token.IMPORT,
			},
			expectedInterface: nil,
			expectedStruct:    nil,
		},
		{
			name:          "InterfaceGenDecl",
			depth:         2,
			origPkgName:   "controller",
			interfaceName: "Controller",
			node:          interfaceGenDecl,
			expectedNode:  interfaceGenDecl,
			expectedInterface: &interfaceType{
				Exported: true,
				Name:     "Controller",
			},
			expectedStruct: nil,
		},
		{
			name:              "StructGenDecl",
			depth:             2,
			origPkgName:       "controller",
			interfaceName:     "Controller",
			node:              structGenDecl,
			expectedNode:      structGenDecl,
			expectedInterface: nil,
			expectedStruct: &structType{
				Exported: false,
				Name:     "controller",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &genericTypeModifier{
				modifier: modifier{
					depth:  tc.depth,
					logger: clogger,
				},
			}

			node := m.Modify(tc.origPkgName, tc.interfaceName, tc.node)

			assert.Equal(t, tc.expectedNode, node)
			assert.Equal(t, tc.expectedInterface, m.outputs.Interface)
			assert.Equal(t, tc.expectedStruct, m.outputs.Struct)
		})
	}
}

func TestGenericFuncModifier(t *testing.T) {
	logger := log.New(log.None)
	clogger := &log.ColorfulLogger{
		Red:     logger,
		Green:   logger,
		Yellow:  logger,
		Blue:    logger,
		Magenta: logger,
		Cyan:    logger,
		White:   logger,
	}

	exportedFunc := &ast.FuncDecl{
		Recv: nil,
		Name: &ast.Ident{Name: "NewController"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "userGateway"},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "gateway"},
								Sel: &ast.Ident{Name: "UserGateway"},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{Name: "Controller"},
					},
					{
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							X: &ast.CompositeLit{
								Type: &ast.Ident{Name: "controller"},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key:   &ast.Ident{Name: "userGateway"},
										Value: &ast.Ident{Name: "userGateway"},
									},
								},
							},
						},
						&ast.Ident{Name: "nil"},
					},
				},
			},
		},
	}

	unexportedFunc := &ast.FuncDecl{
		Recv: nil,
		Name: &ast.Ident{Name: "newController"},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "ug"},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "gateway"},
								Sel: &ast.Ident{Name: "UserGateway"},
							},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.Ident{Name: "controller"},
						},
					},
					{
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
	}

	exportedMethod := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "c"},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "controller"},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Calculate",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "a"},
							{Name: "b"},
						},
						Type: &ast.StarExpr{
							X: &ast.Ident{Name: "int"},
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "entity"},
								Sel: &ast.Ident{Name: "CalculateResponse"},
							},
						},
					},
					{
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.UnaryExpr{
							X: &ast.CompositeLit{
								Type: &ast.StarExpr{
									X: &ast.SelectorExpr{
										X:   &ast.Ident{Name: "entity"},
										Sel: &ast.Ident{Name: "CalculateResponse"},
									},
								},
								Elts: []ast.Expr{
									&ast.KeyValueExpr{
										Key: &ast.Ident{Name: "Result"},
										Value: &ast.BinaryExpr{
											Op: token.Lookup("*"),
											X:  &ast.Ident{Name: "a"},
											Y:  &ast.Ident{Name: "b"},
										},
									},
								},
							},
						},
						&ast.Ident{Name: "nil"},
					},
				},
			},
		},
	}

	unexportedMethod := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "c"},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "controller"},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "calculate",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "a"},
							{Name: "b"},
						},
						Type: &ast.Ident{Name: "int"},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "resp"},
						},
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "entity"},
							Sel: &ast.Ident{Name: "CalculateResponse"},
						},
					},
					{
						Names: []*ast.Ident{
							{Name: "err"},
						},
						Type: &ast.Ident{Name: "error"},
					},
				},
			},
		},
	}

	voidMethod := &ast.FuncDecl{
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{
						{Name: "h"},
					},
					Type: &ast.StarExpr{
						X: &ast.Ident{Name: "handler"},
					},
				},
			},
		},
		Name: &ast.Ident{
			Name: "Calculate",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{Name: "w"},
						},
						Type: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "http"},
							Sel: &ast.Ident{Name: "ResponseWriter"},
						},
					},
					{
						Names: []*ast.Ident{
							{Name: "r"},
						},
						Type: &ast.StarExpr{
							X: &ast.SelectorExpr{
								X:   &ast.Ident{Name: "http"},
								Sel: &ast.Ident{Name: "Request"},
							},
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "w"},
							Sel: &ast.Ident{Name: "WriteHeader"},
						},
						Args: []ast.Expr{
							&ast.BasicLit{
								Value: "200",
							},
						},
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		depth         int
		origPkgName   string
		interfaceName string
		structName    string
		node          ast.Node
		expectedNode  ast.Node
		expectedFunc  funcType
	}{
		{
			name:          "ExportedFunc",
			depth:         2,
			origPkgName:   "_controller",
			interfaceName: "Controller",
			structName:    "controller",
			node:          exportedFunc,
			expectedNode:  exportedFunc,
			expectedFunc: funcType{
				Exported: true,
				Name:     "NewController",
				Receiver: nil,
				Inputs: fields{
					{
						Names:   []string{"userGateway"},
						Star:    true,
						Package: "gateway",
						Type:    "UserGateway",
					},
				},
				Outputs: fields{
					{
						Names:   nil,
						Star:    false,
						Package: "_controller",
						Type:    "Controller",
					},
					{
						Names:   nil,
						Star:    false,
						Package: "",
						Type:    "error",
					},
				},
			},
		},
		{
			name:         "UnexportedFunc",
			depth:        2,
			node:         unexportedFunc,
			expectedNode: unexportedFunc,
			expectedFunc: funcType{
				Exported: false,
				Name:     "newController",
				Receiver: nil,
				Inputs: fields{
					{
						Names:   []string{"ug"},
						Star:    true,
						Package: "gateway",
						Type:    "UserGateway",
					},
				},
				Outputs: fields{
					{
						Names:   nil,
						Star:    true,
						Package: "",
						Type:    "controller",
					},
					{
						Names:   nil,
						Star:    false,
						Package: "",
						Type:    "error",
					},
				},
			},
		},
		{
			name:         "ExportedMethod",
			depth:        2,
			node:         exportedMethod,
			expectedNode: exportedMethod,
			expectedFunc: funcType{
				Exported: true,
				Name:     "Calculate",
				Receiver: &receiver{
					Name: "c",
					Star: true,
					Type: "controller",
				},
				Inputs: fields{
					{
						Names:   []string{"a", "b"},
						Star:    true,
						Package: "",
						Type:    "int",
					},
				},
				Outputs: fields{
					{
						Names:   nil,
						Star:    true,
						Package: "entity",
						Type:    "CalculateResponse",
					},
					{
						Names:   nil,
						Star:    false,
						Package: "",
						Type:    "error",
					},
				},
			},
		},
		{
			name:         "UnexportedMethod",
			depth:        2,
			node:         unexportedMethod,
			expectedNode: unexportedMethod,
			expectedFunc: funcType{
				Exported: false,
				Name:     "calculate",
				Receiver: &receiver{
					Name: "c",
					Star: true,
					Type: "controller",
				},
				Inputs: fields{
					{
						Names:   []string{"a", "b"},
						Star:    false,
						Package: "",
						Type:    "int",
					},
				},
				Outputs: fields{
					{
						Names:   []string{"resp"},
						Star:    false,
						Package: "entity",
						Type:    "CalculateResponse",
					},
					{
						Names:   []string{"err"},
						Star:    false,
						Package: "",
						Type:    "error",
					},
				},
			},
		},
		{
			name:         "VoidMethod",
			depth:        2,
			node:         voidMethod,
			expectedNode: voidMethod,
			expectedFunc: funcType{
				Exported: true,
				Name:     "Calculate",
				Receiver: &receiver{
					Name: "h",
					Star: true,
					Type: "handler",
				},
				Inputs: fields{
					{
						Names:   []string{"w"},
						Star:    false,
						Package: "http",
						Type:    "ResponseWriter",
					},
					{
						Names:   []string{"r"},
						Star:    true,
						Package: "http",
						Type:    "Request",
					},
				},
				Outputs: nil,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			m := &genericFuncModifier{
				modifier: modifier{
					depth:  tc.depth,
					logger: clogger,
				},
			}

			node := m.Modify(tc.origPkgName, tc.interfaceName, tc.structName, tc.node)

			assert.Equal(t, tc.expectedNode, node)
			assert.Equal(t, tc.expectedFunc, m.outputs.Func)
		})
	}
}
