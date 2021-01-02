package modifier

import (
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/moorara/gelato/internal/log"
)

const (
	handlerPkg    = "handler"
	controllerPkg = "controller"
	gatewayPkg    = "gateway"
	repositoryPkg = "repository"
)

func getDecoratedPkgName(origPkgName string) string {
	return "_" + origPkgName
}

func isGenericPkgName(pkgName string) bool {
	return pkgName == handlerPkg ||
		pkgName == controllerPkg ||
		pkgName == gatewayPkg ||
		pkgName == repositoryPkg
}

func isGenericPkgPath(pkgPath string) bool {
	return strings.HasSuffix(pkgPath, "/"+handlerPkg) || strings.Contains(pkgPath, "/"+handlerPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+controllerPkg) || strings.Contains(pkgPath, "/"+controllerPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+gatewayPkg) || strings.Contains(pkgPath, "/"+gatewayPkg+"/") ||
		strings.HasSuffix(pkgPath, "/"+repositoryPkg) || strings.Contains(pkgPath, "/"+repositoryPkg+"/")
}

// MainModifier is used for decorating the main package.
// It implements the Pre and Post astutil.ApplyFunc functions.
type MainModifier struct {
	modifier
	isCallExpr     bool
	isSelectorExpr bool
	importSpecs    []*ast.ImportSpec
	inputs         struct {
		module string
		decDir string
	}
	outputs struct{}
}

// NewMain creates a new main modifier.
func NewMain(depth int, logger *log.ColorfulLogger) *MainModifier {
	m := modifier{
		depth:  depth,
		logger: logger,
	}

	return &MainModifier{
		modifier: m,
	}
}

// Modify modifies an ast.File node.
func (m *MainModifier) Modify(module, decDir string, n ast.Node) ast.Node {
	m.isCallExpr = false
	m.isSelectorExpr = false
	m.importSpecs = nil
	m.inputs.module = module
	m.inputs.decDir = decDir

	return astutil.Apply(n, m.pre, m.post)
}

// Pre is called for each node before the node's children are traversed (pre-order).
func (m *MainModifier) pre(c *astutil.Cursor) bool {
	m.depth++

	switch n := c.Node().(type) {
	case *ast.ImportSpec:
		if pkgPath := strings.Trim(n.Path.Value, `"`); isGenericPkgPath(pkgPath) {
			var pkgName string
			if n.Name != nil {
				pkgName = n.Name.Name
			}
			if pkgName == "" {
				pkgName = filepath.Base(pkgPath)
			}

			new := filepath.Join(m.inputs.module, m.inputs.decDir)
			decPkgPath := strings.Replace(pkgPath, m.inputs.module, new, 1)

			m.importSpecs = append(m.importSpecs, &ast.ImportSpec{
				Name: &ast.Ident{
					// TODO: Resolve NamePos
					Name: getDecoratedPkgName(pkgName),
				},
				Path: &ast.BasicLit{
					// TODO: Resolve ValuePos
					Value: fmt.Sprintf("%q", decPkgPath),
				},
			})
		}

	case *ast.CallExpr:
		m.isCallExpr = true

	case *ast.SelectorExpr:
		m.isSelectorExpr = true

	case *ast.Ident:
		if m.isCallExpr && m.isSelectorExpr && c.Name() == "X" {
			if isGenericPkgName(n.Name) {
				n.Name = getDecoratedPkgName(n.Name)
			}
		}
	}

	return true
}

// Post is called for each node after its children are traversed (post-order).
func (m *MainModifier) post(c *astutil.Cursor) bool {
	m.depth--

	switch n := c.Node().(type) {
	case *ast.GenDecl:
		switch n.Tok {
		case token.IMPORT:
			for _, s := range m.importSpecs {
				n.Specs = append(n.Specs, s)
			}
		}

	case *ast.CallExpr:
		m.isCallExpr = false

	case *ast.SelectorExpr:
		m.isSelectorExpr = false
	}

	return true
}
