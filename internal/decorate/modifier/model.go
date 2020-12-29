package modifier

import "go/ast"

const (
	errorID          = "err"
	implementationID = "impl"
)

type receiver struct {
	Name string
	Star bool
	Type string
}

type field struct {
	Names   []string
	Star    bool
	Package string
	Type    string
}

type fields []field

func (f *fields) Append(n *ast.Field) {
	new := field{}
	for _, id := range n.Names {
		new.Names = append(new.Names, id.Name)
	}
	*f = append(*f, new)
}

func (f *fields) SetStar() {
	i := len(*f) - 1
	(*f)[i].Star = true
}

func (f *fields) SetPackage(n *ast.Ident) {
	i := len(*f) - 1
	(*f)[i].Package = n.Name
}

func (f *fields) SetType(n *ast.Ident) {
	i := len(*f) - 1
	(*f)[i].Type = n.Name
}

type funcType struct {
	Exported bool
	Name     string
	Receiver *receiver
	Inputs   fields
	Outputs  fields
}

type interfaceType struct {
	Exported bool
	Name     string
}

type structType struct {
	Exported bool
	Name     string
}
